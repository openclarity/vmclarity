// Copyright © 2022 Cisco Systems, Inc. and its affiliates.
// All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package scanner

import (
	"fmt"
	"sync/atomic"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/openclarity/vmclarity/runtime_scan/pkg/config"
	types2 "github.com/openclarity/vmclarity/runtime_scan/pkg/types"
)

// run jobs.
func (s *Scanner) jobBatchManagement(scanDone chan struct{}) {
	s.Lock()
	imageIDToScanData := s.instanceIDToScanData
	numberOfWorkers := s.scanConfig.MaxScanParallelism
	imagesStartedToScan := &s.progress.ImagesStartedToScan
	imagesCompletedToScan := &s.progress.ImagesCompletedToScan
	s.Unlock()

	// queue of scan data
	q := make(chan *scanData)
	// done channel takes the result of the job
	done := make(chan bool)

	fullScanDone := make(chan bool)

	// spawn workers
	for i := 0; i < numberOfWorkers; i++ {
		go s.worker(q, i, done, s.killSignal)
	}

	// wait until scan of all images is done - non blocking. once all done, notify on fullScanDone chan
	go func() {
		for c := 0; c < len(imageIDToScanData); c++ {
			select {
			case <-done:
				atomic.AddUint32(imagesCompletedToScan, 1)
			case <-s.killSignal:
				log.WithFields(s.logFields).Debugf("Scan process was canceled - stop waiting for finished jobs")
				return
			}
		}

		fullScanDone <- true
	}()

	// send all scan data on scan data queue, for workers to pick it up.
	for _, data := range imageIDToScanData {
		go func(data *scanData, ks chan bool) {
			select {
			case q <- data:
				atomic.AddUint32(imagesStartedToScan, 1)
			case <-ks:
				log.WithFields(s.logFields).Debugf("Scan process was canceled. imageID=%v, scanUUID=%v", data.instance.ID, data.scanUUID)
				return
			}
		}(data, s.killSignal)
	}

	// wait for killSignal or fullScanDone
	select {
	case <-s.killSignal:
		log.WithFields(s.logFields).Info("Scan process was canceled")
	case <-fullScanDone:
		log.WithFields(s.logFields).Infof("All jobs has finished")
		// Nonblocking notification of a finished scan
		nonBlockingNotification(scanDone)
	}
}

// worker waits for data on the queue, runs a scan job and waits for results from that scan job. Upon completion, done is notified to the caller.
func (s *Scanner) worker(queue chan *scanData, workNumber int, done, ks chan bool) {
	for {
		select {
		case data := <-queue:
			job, err := s.runJob(data)
			if err != nil {
				errMsg := fmt.Errorf("failed to run job: %v", err)
				log.WithFields(s.logFields).Error(errMsg)
				s.Lock()
				data.success = false
				data.scanErr = &types2.ScanError{
					ErrMsg:    err.Error(),
					ErrType:   string(types2.JobRun),
					ErrSource: types2.ScanErrSourceJob,
				}
				data.completed = true
				s.Unlock()
			} else {
				s.waitForResult(data, ks)
			}

			s.deleteJobIfNeeded(&job, data.success, data.completed)

			select {
			case done <- true:
			case <-ks:
				log.WithFields(s.logFields).Infof("Image scan was canceled. imageID=%v", data.instance.ID)
			}
		case <-ks:
			log.WithFields(s.logFields).Debugf("worker #%v halted", workNumber)
			return
		}
	}
}

func (s *Scanner) waitForResult(data *scanData, ks chan bool) {
	//log.WithFields(s.logFields).Infof("Waiting for result. imageID=%+v", data.imageID)
	ticker := time.NewTicker(s.scanConfig.JobResultTimeout)
	select {
	case <-data.resultChan:
		log.WithFields(s.logFields).Infof("Instance scanned result has arrived. instanceID=%v", data.instance.ID)
	case <-ticker.C:
		errMsg := fmt.Errorf("job has timed out. imageID=%v", data.instance.ID)
		log.WithFields(s.logFields).Warn(errMsg)
		s.Lock()
		data.success = false
		data.scanErr = &types2.ScanError{
			ErrMsg:    errMsg.Error(),
			ErrType:   string(types2.JobTimeout),
			ErrSource: types2.ScanErrSourceJob,
		}
		data.timeout = true
		data.completed = true
		s.Unlock()
	case <-ks:
		log.WithFields(s.logFields).Infof("Image scan was canceled. imageID=%v", data.instance.ID)
	}
}

func (s *Scanner) runJob(data *scanData) (types2.Job, error) {
	rootVolume, err := s.providerClient.GetInstanceRootVolume(data.instance)
	if err != nil {
		return types2.Job{}, fmt.Errorf("failed to get instance root volume. instance id=%v: %v", data.instance.ID, err)
	}

	// create a snapshot of that vm
	srcSnapshot, err := s.providerClient.CreateSnapshot(rootVolume)
	if err != nil {
		return types2.Job{}, fmt.Errorf("failed to create snapshot: %v", err)
	}
	if err := s.providerClient.WaitForSnapshotReady(srcSnapshot); err != nil {
		return types2.Job{}, err
	}

	//copy the snapshot to the scanner region
	// TODO make sure we need this.
	cpySnapshot, err := s.providerClient.CopySnapshot(srcSnapshot, s.region)
	if err != nil {
		return types2.Job{}, fmt.Errorf("failed to copy snapshot: %v", err)
	}
	if err := s.providerClient.WaitForSnapshotReady(cpySnapshot); err != nil {
		return types2.Job{}, err
	}

	// create the scanner job (vm) with a boot script
	launchedInstance, err := s.providerClient.LaunchInstance(s.amiID, "", cpySnapshot)
	if err != nil {
		return types2.Job{}, fmt.Errorf("failed to launch instance: %v", err)
	}

	return types2.Job{
		Instance:    launchedInstance,
		SrcSnapshot: srcSnapshot,
		DstSnapshot: cpySnapshot,
	}, nil
}

func (s *Scanner) deleteJobIfNeeded(job *types2.Job, isSuccessfulJob, isCompletedJob bool) {
	if job == nil {
		return
	}

	// delete uncompleted jobs - scan process was canceled
	if !isCompletedJob {
		s.deleteJob(job)
		return
	}

	switch s.scanConfig.DeleteJobPolicy {
	case config.DeleteJobPolicyNever:
		// do nothing
	case config.DeleteJobPolicyAll:
		s.deleteJob(job)
	case config.DeleteJobPolicySuccessful:
		if isSuccessfulJob {
			s.deleteJob(job)
		}
	}
}

func (s *Scanner) deleteJob(job *types2.Job) {
	if err := s.providerClient.DeleteInstance(job.Instance); err != nil {
		log.Errorf("failed to delete instance: %v", err)
	}
	if err := s.providerClient.DeleteSnapshot(job.SrcSnapshot); err != nil {
		log.Errorf("failed to delete source snapshot: %v", err)
	}
	if err := s.providerClient.DeleteSnapshot(job.DstSnapshot); err != nil {
		log.Errorf("failed to delete dest snapshot: %v", err)
	}
}

// Due to K8s names constraint we will take the image name w/o the tag and repo.
//func getSimpleImageName(imageName string) (string, error) {
//	ref, err := reference.ParseNormalizedNamed(imageName)
//	if err != nil {
//		return "", fmt.Errorf("failed to parse image name. name=%v: %v", imageName, err)
//	}
//
//	refName := ref.Name()
//	// Take only image name from repo path (ex. solsson/kafka ==> kafka)
//	repoEnd := strings.LastIndex(refName, "/")
//
//	return refName[repoEnd+1:], nil
//}

// Job names require their names to follow the DNS label standard as defined in RFC 1123
// Note: job name is added as a label to the pod template spec so it should follow the DNS label standard and not just DNS-1123 subdomain
//
// This means the name must:
// * contain at most 63 characters
// * contain only lowercase alphanumeric characters or ‘-’
// * start with an alphanumeric character
// * end with an alphanumeric character.
//func createJobName(imageName string) (string, error) {
//	//simpleName, err := getSimpleImageName(imageName)
//	//if err != nil {
//	//	return "", err
//	//}
//
//	jobName := "scanner-" + simpleName + "-" + uuid.NewV4().String()
//
//	// contain at most 63 characters
//	jobName = stringsutils.TruncateString(jobName, k8s.MaxK8sJobName)
//
//	// contain only lowercase alphanumeric characters or ‘-’
//	jobName = strings.ToLower(jobName)
//	jobName = strings.ReplaceAll(jobName, "_", "-")
//
//	// no need to validate start, we are using 'jobName'
//
//	// end with an alphanumeric character
//	jobName = strings.TrimRight(jobName, "-")
//
//	return jobName, nil
//}

//func (s *Scanner) createJob(data *scanData) (*batchv1.Job, error) {
//	// We will scan each image once, based on the first pod context. The result will be applied for all other pods with this image.
//	podContext := data.contexts[0]
//
//	jobName, err := createJobName(podContext.imageName)
//	if err != nil {
//		return nil, fmt.Errorf("failed to create job name. namespace=%v, pod=%v, container=%v, image=%v, hash=%v: %v",
//			podContext.namespace, podContext.podName, podContext.containerName, podContext.imageName, data.imageHash, err)
//	}
//
//	// Set job values on scanner job template
//	job := s.scannerJobTemplate.DeepCopy()
//	if !data.shouldScanCISDockerBenchmark {
//		removeCISDockerBenchmarkScannerFromJob(job)
//	}
//	job.SetName(jobName)
//	job.SetNamespace(podContext.namespace)
//	setJobScanUUID(job, data.scanUUID)
//	setJobImageIDToScan(job, data.imageID)
//	setJobImageHashToScan(job, data.imageHash)
//	setJobImageNameToScan(job, podContext.imageName)
//	if podContext.imagePullSecret != "" {
//		log.WithFields(s.logFields).Debugf("Adding private registry credentials to image: %s", podContext.imageName)
//		setJobDockerConfigFromImagePullSecret(job, podContext.imagePullSecret)
//	} else {
//		// Use private repo sa credentials only if there is no imagePullSecret
//		for _, adder := range s.credentialAdders {
//			if adder.ShouldAdd() {
//				adder.Add(job)
//			}
//		}
//	}
//
//	return job, nil
//}
//
//func removeCISDockerBenchmarkScannerFromJob(job *batchv1.Job) {
//	var containers []corev1.Container
//	for i := range job.Spec.Template.Spec.Containers {
//		container := job.Spec.Template.Spec.Containers[i]
//		if container.Name != cisDockerBenchmarkScannerContainerName {
//			containers = append(containers, container)
//		}
//	}
//	job.Spec.Template.Spec.Containers = containers
//}

// Create docker config from imagePullSecret that contains the username and the password required to pull the image.
// We need to do the following:
// 1. Create a volume that holds the `secretName` data
// 2. Mount the volume into each container to a specific path (`BasicVolumeMountPath`/`DockerConfigFileName`)
// 3. Set `DOCKER_CONFIG` to point to the directory that contains the config.json.
//func setJobDockerConfigFromImagePullSecret(job *batchv1.Job, secretName string) {
//	job.Spec.Template.Spec.Volumes = append(job.Spec.Template.Spec.Volumes, corev1.Volume{
//		Name: _creds.BasicVolumeName,
//		VolumeSource: corev1.VolumeSource{
//			Secret: &corev1.SecretVolumeSource{
//				SecretName: secretName,
//				Items: []corev1.KeyToPath{
//					{
//						Key:  corev1.DockerConfigJsonKey,
//						Path: _creds.DockerConfigFileName,
//					},
//				},
//			},
//		},
//	})
//	for i := range job.Spec.Template.Spec.Containers {
//		container := &job.Spec.Template.Spec.Containers[i]
//		container.VolumeMounts = append(container.VolumeMounts, corev1.VolumeMount{
//			Name:      _creds.BasicVolumeName,
//			ReadOnly:  true,
//			MountPath: _creds.BasicVolumeMountPath,
//		})
//		container.Env = append(container.Env, corev1.EnvVar{
//			Name:  _creds.DockerConfigEnvVar,
//			Value: _creds.BasicVolumeMountPath,
//		})
//	}
//}

//func setJobImageIDToScan(job *batchv1.Job, imageID string) {
//	for i := range job.Spec.Template.Spec.Containers {
//		container := &job.Spec.Template.Spec.Containers[i]
//		container.Env = append(container.Env, corev1.EnvVar{Name: shared.ImageIDToScan, Value: imageID})
//	}
//}
//
//func setJobImageHashToScan(job *batchv1.Job, imageHash string) {
//	for i := range job.Spec.Template.Spec.Containers {
//		container := &job.Spec.Template.Spec.Containers[i]
//		container.Env = append(container.Env, corev1.EnvVar{Name: shared.ImageHashToScan, Value: imageHash})
//	}
//}
//
//func setJobImageNameToScan(job *batchv1.Job, imageName string) {
//	for i := range job.Spec.Template.Spec.Containers {
//		container := &job.Spec.Template.Spec.Containers[i]
//		container.Env = append(container.Env, corev1.EnvVar{Name: shared.ImageNameToScan, Value: imageName})
//	}
//}
//
//func setJobScanUUID(job *batchv1.Job, scanUUID string) {
//	for i := range job.Spec.Template.Spec.Containers {
//		container := &job.Spec.Template.Spec.Containers[i]
//		container.Env = append(container.Env, corev1.EnvVar{Name: shared.ScanUUID, Value: scanUUID})
//	}
//}

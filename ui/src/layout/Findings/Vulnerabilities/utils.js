import { sortBy } from 'lodash';
import CVSS from '@turingpointde/cvss.js';

export const getHigestVersionCvssData = (cvssData) => {
    const sortedCvss = sortBy(cvssData || [], ["version"]);

    const {vector, metrics} = sortedCvss[0];
    const cvssVector = CVSS(vector);

    return {
        score: cvssVector.getScore(),
        temporalScore: cvssVector.getTemporalScore(),
        environmentalScore: cvssVector.getEnvironmentalScore(),
        severity: cvssVector.getRating(),
        metrics: cvssVector.getDetailedVectorObject().metrics,
        vector,
        exploitabilityScore: metrics.exploitabilityScore,
        impactScore: metrics.impactScore
    }
}
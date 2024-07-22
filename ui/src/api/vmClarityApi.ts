import { VMClarityApi } from './generated';
import axiosClient from './axiosClient';

const vmClarityApi = new VMClarityApi(undefined, undefined, axiosClient);

export default vmClarityApi;

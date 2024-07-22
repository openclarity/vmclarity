import axios from 'axios';

const axiosClient = axios.create({
  baseURL: `${window.location.origin}/api`,
});

export default axiosClient;

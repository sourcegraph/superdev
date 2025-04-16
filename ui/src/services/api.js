import axios from 'axios';

const API_BASE_URL = 'http://localhost:8080';

export const fetchThreads = async () => {
  const response = await axios.get(`${API_BASE_URL}/threads`);
  return response.data;
};

export const fetchThreadOutput = async (threadId) => {
  const response = await axios.get(`${API_BASE_URL}/output?thread_id=${threadId}`);
  return response.data;
};
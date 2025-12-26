import api from './api';

export const postService = {
  async getPosts(page = 1, pageSize = 20, sort = 'new') {
    const response = await api.get('/posts', {
      params: {
        page,
        page_size: pageSize,
        sort,
      },
    });
    return response.data;
  },

  async getPost(id) {
    const response = await api.get(`/posts/${id}`);
    return response.data;
  },

  async getPostById(id) {
    const response = await api.get(`/posts/${id}`);
    return response.data;
  },

  async createPost(postData) {
    const response = await api.post('/posts', postData);
    return response.data;
  },

  async updatePost(id, postData) {
    const response = await api.put(`/posts/${id}`, postData);
    return response.data;
  },

  async deletePost(id) {
    const response = await api.delete(`/posts/${id}`);
    return response.data;
  },

  async searchPosts(query, page = 1, pageSize = 20) {
    const response = await api.get('/posts/search', {
      params: {
        q: query,
        page,
        page_size: pageSize,
      },
    });
    return response.data;
  },
};

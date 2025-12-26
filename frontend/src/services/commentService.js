import api from './api';

export const commentService = {
  async getCommentsByPost(postId) {
    const response = await api.get(`/posts/${postId}/comments`);
    return response.data;
  },

  async createComment(postId, text, parentId = null) {
    const response = await api.post(`/posts/${postId}/comments`, {
      text,
      parent_id: parentId,
    });
    return response.data;
  },

  async updateComment(commentId, text) {
    const response = await api.put(`/comments/${commentId}`, {
      text,
    });
    return response.data;
  },

  async deleteComment(commentId) {
    const response = await api.delete(`/comments/${commentId}`);
    return response.data;
  },
};

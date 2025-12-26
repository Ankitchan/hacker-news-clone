import api from './api';

export const voteService = {
  async voteOnPost(postId, voteType) {
    // voteType: 1 for upvote, -1 for downvote
    const response = await api.post(`/votes/posts/${postId}`, {
      vote_type: voteType,
    });
    return response.data;
  },

  async voteOnComment(commentId, voteType) {
    // voteType: 1 for upvote, -1 for downvote
    const response = await api.post(`/votes/comments/${commentId}`, {
      vote_type: voteType,
    });
    return response.data;
  },

  async removeVoteOnPost(postId) {
    const response = await api.delete(`/votes/posts/${postId}`);
    return response.data;
  },

  async removeVoteOnComment(commentId) {
    const response = await api.delete(`/votes/comments/${commentId}`);
    return response.data;
  },
};

import { useState, useEffect } from 'react';
import { useAuth } from '../context/AuthContext';
import { commentService } from '../services/commentService';
import Comment from './Comment';
import './CommentList.css';

const CommentList = ({ postId }) => {
  const { isAuthenticated } = useAuth();
  const [comments, setComments] = useState([]);
  const [newCommentText, setNewCommentText] = useState('');
  const [loading, setLoading] = useState(true);
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState('');

  const loadComments = async () => {
    setLoading(true);
    setError('');
    try {
      const data = await commentService.getCommentsByPost(postId);
      setComments(data.comments || []);
    } catch (err) {
      setError('Failed to load comments');
      console.error('Error loading comments:', err);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadComments();
  }, [postId]);

  const handleSubmitComment = async (e) => {
    e.preventDefault();
    if (!newCommentText.trim()) return;

    setSubmitting(true);
    try {
      await commentService.createComment(postId, newCommentText.trim(), null);
      setNewCommentText('');
      await loadComments();
    } catch (err) {
      console.error('Failed to post comment:', err);
      alert('Failed to post comment');
    } finally {
      setSubmitting(false);
    }
  };

  const handleCommentUpdate = () => {
    loadComments();
  };

  const handleCommentDelete = () => {
    loadComments();
  };

  if (loading) {
    return <div className="comments-loading">Loading comments...</div>;
  }

  if (error) {
    return <div className="comments-error">{error}</div>;
  }

  return (
    <div className="comments-section">
      <h2 className="comments-header">
        Comments ({comments.length})
      </h2>

      {isAuthenticated ? (
        <form onSubmit={handleSubmitComment} className="comment-form">
          <textarea
            value={newCommentText}
            onChange={(e) => setNewCommentText(e.target.value)}
            placeholder="Add a comment..."
            className="comment-textarea"
            rows="4"
          />
          <button
            type="submit"
            disabled={submitting || !newCommentText.trim()}
            className="comment-submit-btn"
          >
            {submitting ? 'Posting...' : 'Add Comment'}
          </button>
        </form>
      ) : (
        <div className="comment-login-prompt">
          Please log in to add a comment.
        </div>
      )}

      <div className="comments-list">
        {comments.length === 0 ? (
          <p className="no-comments">No comments yet. Be the first to comment!</p>
        ) : (
          comments.map((comment) => (
            <Comment
              key={comment.id}
              comment={comment}
              onUpdate={handleCommentUpdate}
              onDelete={handleCommentDelete}
              level={0}
            />
          ))
        )}
      </div>
    </div>
  );
};

export default CommentList;

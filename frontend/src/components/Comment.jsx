import { useState } from 'react';
import { Link } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import { commentService } from '../services/commentService';
import { voteService } from '../services/voteService';
import './Comment.css';

const Comment = ({ comment, onUpdate, onDelete, level = 0 }) => {
  const { user, isAuthenticated } = useAuth();
  const [isEditing, setIsEditing] = useState(false);
  const [editText, setEditText] = useState(comment.text);
  const [isReplying, setIsReplying] = useState(false);
  const [replyText, setReplyText] = useState('');
  const [localTotalPoints, setLocalTotalPoints] = useState(comment.total_points || 0);
  const [userVote, setUserVote] = useState(null);
  const [voting, setVoting] = useState(false);
  const [submitting, setSubmitting] = useState(false);

  const formatTimeAgo = (timestamp) => {
    const now = new Date();
    const posted = new Date(timestamp);
    const diffInSeconds = Math.floor((now - posted) / 1000);

    if (diffInSeconds < 60) return 'just now';
    if (diffInSeconds < 3600) return `${Math.floor(diffInSeconds / 60)} minutes ago`;
    if (diffInSeconds < 86400) return `${Math.floor(diffInSeconds / 3600)} hours ago`;
    return `${Math.floor(diffInSeconds / 86400)} days ago`;
  };

  const isOwnComment = user && comment.user_id === user.id;

  const handleEdit = async () => {
    if (!editText.trim()) return;

    setSubmitting(true);
    try {
      await commentService.updateComment(comment.id, editText.trim());
      setIsEditing(false);
      if (onUpdate) onUpdate();
    } catch (error) {
      console.error('Failed to update comment:', error);
      alert('Failed to update comment');
    } finally {
      setSubmitting(false);
    }
  };

  const handleDelete = async () => {
    if (!window.confirm('Are you sure you want to delete this comment?')) return;

    try {
      await commentService.deleteComment(comment.id);
      if (onDelete) onDelete(comment.id);
    } catch (error) {
      console.error('Failed to delete comment:', error);
      alert('Failed to delete comment');
    }
  };

  const handleReply = async () => {
    if (!replyText.trim()) return;

    setSubmitting(true);
    try {
      await commentService.createComment(comment.post_id, replyText.trim(), comment.id);
      setReplyText('');
      setIsReplying(false);
      if (onUpdate) onUpdate();
    } catch (error) {
      console.error('Failed to post reply:', error);
      alert('Failed to post reply');
    } finally {
      setSubmitting(false);
    }
  };

  const handleVote = async (voteType) => {
    if (!isAuthenticated || voting) return;

    setVoting(true);
    try {
      if (userVote === voteType) {
        // Remove vote
        await voteService.removeVoteOnComment(comment.id);
        setUserVote(null);
      } else {
        // Create or change vote
        await voteService.voteOnComment(comment.id, voteType);
        setUserVote(voteType);
      }

      // Trigger reload of entire comment tree to get accurate cumulative counts
      if (onUpdate) {
        onUpdate();
      }
    } catch (error) {
      console.error('Vote failed:', error);
    } finally {
      setVoting(false);
    }
  };

  const maxNestLevel = 5;
  const showReplies = level < maxNestLevel;

  return (
    <div className={`comment ${level > 0 ? 'nested' : ''}`} style={{ marginLeft: `${level * 20}px` }}>
      <div className="comment-header">
        <div className="comment-meta">
          <Link to={`/user/${comment.username}`} className="comment-author">
            {comment.username}
          </Link>
          <span className="comment-separator">|</span>
          <span className="comment-time">{formatTimeAgo(comment.created_at)}</span>
          {comment.updated_at !== comment.created_at && (
            <span className="comment-edited">(edited)</span>
          )}
        </div>
      </div>

      <div className="comment-body">
        {isEditing ? (
          <div className="comment-edit-form">
            <textarea
              value={editText}
              onChange={(e) => setEditText(e.target.value)}
              className="comment-edit-textarea"
              rows="4"
            />
            <div className="comment-edit-actions">
              <button onClick={handleEdit} disabled={submitting} className="btn-primary">
                {submitting ? 'Saving...' : 'Save'}
              </button>
              <button onClick={() => setIsEditing(false)} disabled={submitting} className="btn-secondary">
                Cancel
              </button>
            </div>
          </div>
        ) : (
          <p className="comment-text">{comment.text}</p>
        )}
      </div>

      <div className="comment-actions">
        {isAuthenticated && !isOwnComment && (
          <div className="comment-vote-buttons">
            <button
              className={`vote-btn-small upvote ${userVote === 1 ? 'active' : ''}`}
              onClick={() => handleVote(1)}
              disabled={voting}
              title="Upvote"
            >
              ▲
            </button>
            <span className="comment-points">{localTotalPoints}</span>
            <button
              className={`vote-btn-small downvote ${userVote === -1 ? 'active' : ''}`}
              onClick={() => handleVote(-1)}
              disabled={voting}
              title="Downvote"
            >
              ▼
            </button>
          </div>
        )}
        {!isOwnComment && !isAuthenticated && (
          <span className="comment-points">{localTotalPoints} points</span>
        )}
        {isOwnComment && !isEditing && (
          <>
            <button onClick={() => setIsEditing(true)} className="comment-action-link">
              edit
            </button>
            <span className="comment-separator">|</span>
            <button onClick={handleDelete} className="comment-action-link">
              delete
            </button>
          </>
        )}
        {isAuthenticated && !isEditing && showReplies && (
          <>
            {(isOwnComment || !isAuthenticated) && <span className="comment-separator">|</span>}
            <button onClick={() => setIsReplying(!isReplying)} className="comment-action-link">
              reply
            </button>
          </>
        )}
      </div>

      {isReplying && (
        <div className="comment-reply-form">
          <textarea
            value={replyText}
            onChange={(e) => setReplyText(e.target.value)}
            placeholder="Write your reply..."
            className="comment-reply-textarea"
            rows="3"
          />
          <div className="comment-reply-actions">
            <button onClick={handleReply} disabled={submitting || !replyText.trim()} className="btn-primary">
              {submitting ? 'Posting...' : 'Reply'}
            </button>
            <button onClick={() => setIsReplying(false)} disabled={submitting} className="btn-secondary">
              Cancel
            </button>
          </div>
        </div>
      )}

      {comment.replies && comment.replies.length > 0 && showReplies && (
        <div className="comment-replies">
          {comment.replies.map((reply) => (
            <Comment
              key={reply.id}
              comment={reply}
              onUpdate={onUpdate}
              onDelete={onDelete}
              level={level + 1}
            />
          ))}
        </div>
      )}
    </div>
  );
};

export default Comment;

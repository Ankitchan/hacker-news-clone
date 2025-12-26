import { useState } from 'react';
import { Link } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import { voteService } from '../services/voteService';
import './Post.css';

const Post = ({ post, rank, onVoteUpdate }) => {
  const { isAuthenticated, user } = useAuth();
  const [localPoints, setLocalPoints] = useState(post.points);
  const [userVote, setUserVote] = useState(null); // 1 for upvote, -1 for downvote, null for no vote
  const [voting, setVoting] = useState(false);

  const formatTimeAgo = (timestamp) => {
    const now = new Date();
    const posted = new Date(timestamp);
    const diffInSeconds = Math.floor((now - posted) / 1000);

    if (diffInSeconds < 60) return 'just now';
    if (diffInSeconds < 3600) return `${Math.floor(diffInSeconds / 60)} minutes ago`;
    if (diffInSeconds < 86400) return `${Math.floor(diffInSeconds / 3600)} hours ago`;
    return `${Math.floor(diffInSeconds / 86400)} days ago`;
  };

  const getDomain = (url) => {
    if (!url) return '';
    try {
      const urlObj = new URL(url);
      return urlObj.hostname.replace('www.', '');
    } catch {
      return '';
    }
  };

  const handleVote = async (voteType) => {
    if (!isAuthenticated || voting) return;

    setVoting(true);

    try {
      if (userVote === voteType) {
        // Remove vote if clicking same button
        await voteService.removeVoteOnPost(post.id);
        setLocalPoints(prev => prev - voteType);
        setUserVote(null);
      } else {
        // Add or change vote
        await voteService.voteOnPost(post.id, voteType);

        if (userVote === null) {
          // New vote
          setLocalPoints(prev => prev + voteType);
        } else {
          // Changing vote (from upvote to downvote or vice versa)
          setLocalPoints(prev => prev + (voteType * 2));
        }
        setUserVote(voteType);
      }

      // Notify parent component if callback provided
      if (onVoteUpdate) {
        onVoteUpdate(post.id);
      }
    } catch (error) {
      console.error('Vote failed:', error);
    } finally {
      setVoting(false);
    }
  };

  const isUrl = post.url?.Valid && post.url?.String;
  const url = isUrl ? post.url.String : `/post/${post.id}`;
  const isExternal = isUrl;
  const isOwnPost = user && post.user_id === user.id;

  return (
    <div className="post-item">
      <div className="post-vote-section">
        <div className="post-rank">{rank}.</div>
        {isAuthenticated && !isOwnPost && (
          <div className="vote-buttons">
            <button
              className={`vote-btn upvote ${userVote === 1 ? 'active' : ''}`}
              onClick={() => handleVote(1)}
              disabled={voting}
              title="Upvote"
            >
              ▲
            </button>
            <button
              className={`vote-btn downvote ${userVote === -1 ? 'active' : ''}`}
              onClick={() => handleVote(-1)}
              disabled={voting}
              title="Downvote"
            >
              ▼
            </button>
          </div>
        )}
      </div>
      <div className="post-content">
        <div className="post-title-row">
          {isExternal ? (
            <a href={url} target="_blank" rel="noopener noreferrer" className="post-title">
              {post.title}
            </a>
          ) : (
            <Link to={url} className="post-title">
              {post.title}
            </Link>
          )}
          {isUrl && <span className="post-domain">({getDomain(post.url.String)})</span>}
        </div>
        <div className="post-meta">
          <span className="post-points">{localPoints} points</span>
          <span className="post-separator">by</span>
          <Link to={`/user/${post.username}`} className="post-author">
            {post.username}
          </Link>
          <span className="post-separator">|</span>
          <span className="post-time">{formatTimeAgo(post.created_at)}</span>
          <span className="post-separator">|</span>
          <Link to={`/post/${post.id}`} className="post-comments">
            {post.comment_count || 0} comments
          </Link>
        </div>
      </div>
    </div>
  );
};

export default Post;

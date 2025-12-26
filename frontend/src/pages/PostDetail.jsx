import { useState, useEffect } from 'react';
import { useParams, useNavigate, Link } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import { postService } from '../services/postService';
import CommentList from '../components/CommentList';
import './PostDetail.css';

const PostDetail = () => {
  const { id } = useParams();
  const navigate = useNavigate();
  const { user, logout } = useAuth();
  const [post, setPost] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    const loadPost = async () => {
      setLoading(true);
      setError('');
      try {
        const data = await postService.getPostById(id);
        setPost(data.post);
      } catch (err) {
        setError('Failed to load post');
        console.error('Error loading post:', err);
      } finally {
        setLoading(false);
      }
    };

    loadPost();
  }, [id]);

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

  const handleLogout = () => {
    logout();
    navigate('/login');
  };

  if (loading) {
    return (
      <div className="post-detail-container">
        <div className="loading">Loading post...</div>
      </div>
    );
  }

  if (error || !post) {
    return (
      <div className="post-detail-container">
        <div className="error">{error || 'Post not found'}</div>
        <Link to="/" className="back-link">← Back to home</Link>
      </div>
    );
  }

  const isUrl = post.url?.Valid && post.url?.String;
  const url = isUrl ? post.url.String : null;

  return (
    <div className="post-detail-container">
      <header className="post-detail-header">
        <div className="header-content">
          <Link to="/" className="site-title-link">
            <h1 className="site-title">Hacker News</h1>
          </Link>
          <div className="header-right">
            {user ? (
              <>
                <span className="username">{user.username}</span>
                <span className="nav-separator">|</span>
                <button className="nav-link" onClick={handleLogout}>
                  logout
                </button>
              </>
            ) : (
              <button className="nav-link" onClick={() => navigate('/login')}>
                login
              </button>
            )}
          </div>
        </div>
      </header>

      <main className="post-detail-main">
        <div className="post-detail-content">
          <div className="post-detail-header-section">
            <h1 className="post-detail-title">
              {isUrl ? (
                <a href={url} target="_blank" rel="noopener noreferrer">
                  {post.title}
                </a>
              ) : (
                post.title
              )}
            </h1>
            {isUrl && <span className="post-detail-domain">({getDomain(url)})</span>}
          </div>

          <div className="post-detail-meta">
            <span className="post-points">{post.points} points</span>
            <span className="post-separator">by</span>
            <Link to={`/user/${post.username}`} className="post-author">
              {post.username}
            </Link>
            <span className="post-separator">|</span>
            <span className="post-time">{formatTimeAgo(post.created_at)}</span>
          </div>

          {post.text && (
            <div className="post-detail-text">
              <p>{post.text}</p>
            </div>
          )}
        </div>

        <CommentList postId={id} />
      </main>
    </div>
  );
};

export default PostDetail;

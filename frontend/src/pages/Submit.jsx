import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import { postService } from '../services/postService';
import './Submit.css';

const Submit = () => {
  const [title, setTitle] = useState('');
  const [url, setUrl] = useState('');
  const [text, setText] = useState('');
  const [error, setError] = useState('');
  const [isSubmitting, setIsSubmitting] = useState(false);
  const { user } = useAuth();
  const navigate = useNavigate();

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError('');

    // Validate
    if (!title.trim()) {
      setError('Title is required');
      return;
    }

    if (!url.trim() && !text.trim()) {
      setError('Either URL or description is required');
      return;
    }

    setIsSubmitting(true);

    try {
      const postData = {
        title: title.trim(),
      };

      if (url.trim()) {
        postData.url = url.trim();
      }

      if (text.trim()) {
        postData.text = text.trim();
      }

      const response = await postService.createPost(postData);

      // Navigate to the newly created post
      navigate(`/post/${response.id}`);
    } catch (err) {
      setError(err.response?.data?.error || 'Failed to create post. Please try again.');
    } finally {
      setIsSubmitting(false);
    }
  };

  const handleCancel = () => {
    navigate('/');
  };

  return (
    <div className="submit-container">
      <header className="submit-header">
        <div className="header-content">
          <h1 className="site-title" onClick={() => navigate('/')}>
            Hacker News
          </h1>
        </div>
      </header>

      <main className="submit-main">
        <div className="submit-form-container">
          <h2>Submit</h2>

          {error && <div className="error-message">{error}</div>}

          <form onSubmit={handleSubmit} className="submit-form">
            <div className="form-group">
              <label htmlFor="title">Title</label>
              <input
                type="text"
                id="title"
                value={title}
                onChange={(e) => setTitle(e.target.value)}
                placeholder="Enter post title"
                disabled={isSubmitting}
                autoFocus
              />
            </div>

            <div className="form-group">
              <label htmlFor="url">URL</label>
              <input
                type="url"
                id="url"
                value={url}
                onChange={(e) => setUrl(e.target.value)}
                placeholder="https://example.com"
                disabled={isSubmitting}
              />
            </div>

            <div className="form-divider">
              <span>or</span>
            </div>

            <div className="form-group">
              <label htmlFor="text">Description</label>
              <textarea
                id="text"
                value={text}
                onChange={(e) => setText(e.target.value)}
                placeholder="Enter text description"
                rows="6"
                disabled={isSubmitting}
              />
            </div>

            <div className="form-info">
              Leave URL blank to submit a question for discussion. If there is no URL,
              description will be used for the post content.
            </div>

            <div className="form-actions">
              <button
                type="submit"
                className="submit-btn"
                disabled={isSubmitting}
              >
                {isSubmitting ? 'Submitting...' : 'Submit'}
              </button>
              <button
                type="button"
                className="cancel-btn"
                onClick={handleCancel}
                disabled={isSubmitting}
              >
                Cancel
              </button>
            </div>
          </form>
        </div>
      </main>
    </div>
  );
};

export default Submit;

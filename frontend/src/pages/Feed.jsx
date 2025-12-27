import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import PostList from '../components/PostList';
import SearchBar from '../components/SearchBar';
import api from '../services/api';
import './Feed.css';

const Feed = () => {
  const [sort, setSort] = useState('new');
  const [searchQuery, setSearchQuery] = useState('');
  const [unreadCount, setUnreadCount] = useState(0);
  const { user, logout } = useAuth();
  const navigate = useNavigate();

  // Fetch unread notification count
  useEffect(() => {
    if (!user) return;

    const fetchUnreadCount = async () => {
      try {
        const response = await api.get('/notifications/unread/count');
        setUnreadCount(response.data.unread_count || 0);
      } catch (error) {
        console.error('Failed to fetch unread count:', error);
      }
    };

    fetchUnreadCount();
    // Poll for updates every 30 seconds
    const interval = setInterval(fetchUnreadCount, 30000);

    return () => clearInterval(interval);
  }, [user]);

  const handleLogout = () => {
    logout();
    navigate('/login');
  };

  const handleSearch = (query) => {
    setSearchQuery(query);
    setSort('new'); // Reset sort when searching
  };

  const handleClear = () => {
    setSearchQuery('');
  };

  return (
    <div className="feed-container">
      <header className="feed-header">
        <div className="header-content">
          <div className="header-left">
            <h1 className="site-title">Hacker News</h1>
            <nav className="main-nav">
              <button
                className={`nav-link ${sort === 'new' ? 'active' : ''}`}
                onClick={() => setSort('new')}
              >
                new
              </button>
              <span className="nav-separator">|</span>
              <button
                className={`nav-link ${sort === 'top' ? 'active' : ''}`}
                onClick={() => setSort('top')}
              >
                top
              </button>
              <span className="nav-separator">|</span>
              <button
                className={`nav-link ${sort === 'best' ? 'active' : ''}`}
                onClick={() => setSort('best')}
              >
                best
              </button>
              <span className="nav-separator">|</span>
              <button className="nav-link" onClick={() => navigate('/submit')}>
                submit
              </button>
            </nav>
          </div>
          <div className="header-right">
            {user ? (
              <>
                <button className="notification-bell" onClick={() => navigate('/notifications')} title="Notifications">
                  <span className="bell-icon">🔔</span>
                  {unreadCount > 0 && (
                    <span className="notification-badge">{unreadCount}</span>
                  )}
                </button>
                <span className="nav-separator">|</span>
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

      <SearchBar onSearch={handleSearch} onClear={handleClear} currentSearchQuery={searchQuery} />

      <main className="feed-main">
        <PostList key={searchQuery || sort} sort={sort} searchQuery={searchQuery} />
      </main>
    </div>
  );
};

export default Feed;

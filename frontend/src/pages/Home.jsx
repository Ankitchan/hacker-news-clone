import { useAuth } from '../context/AuthContext';
import { useNavigate } from 'react-router-dom';
import './Home.css';

const Home = () => {
  const { user, logout } = useAuth();
  const navigate = useNavigate();

  const handleLogout = () => {
    logout();
    navigate('/login');
  };

  return (
    <div className="home-container">
      <header className="header">
        <div className="header-content">
          <h1>Hacker News Clone</h1>
          {user && (
            <div className="user-info">
              <span>Welcome, {user.username}!</span>
              <button onClick={handleLogout} className="logout-btn">
                Logout
              </button>
            </div>
          )}
        </div>
      </header>

      <main className="main-content">
        <div className="welcome-message">
          <h2>Welcome to Hacker News Clone</h2>
          <p>You are successfully logged in!</p>
          <p className="user-details">
            <strong>Username:</strong> {user?.username}<br />
            <strong>Email:</strong> {user?.email}
          </p>
        </div>
      </main>
    </div>
  );
};

export default Home;

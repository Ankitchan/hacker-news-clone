# Hacker News Clone - Frontend

React frontend for the Hacker News Clone application.

## Features

✅ **Authentication**
- User signup with validation
- User login
- Protected routes
- JWT token management
- Logout functionality

## Tech Stack

- **React** - UI library
- **Vite** - Build tool and dev server
- **React Router** - Client-side routing
- **Axios** - HTTP client
- **Context API** - State management

## Setup

### Prerequisites

- Node.js 25.2.1 or higher
- npm 11.6.4 or higher
- Backend server running on port 8080

### Installation

1. **Install dependencies**
   ```bash
   npm install
   ```

2. **Configure environment**
   ```bash
   # Create .env file (already created)
   # VITE_API_BASE_URL=http://localhost:8080/api
   ```

3. **Start development server**
   ```bash
   npm run dev
   ```

The application will be available at `http://localhost:5173`

## Project Structure

```
frontend/
├── src/
│   ├── components/        # Reusable components
│   │   └── ProtectedRoute.jsx
│   ├── context/          # React Context for state management
│   │   └── AuthContext.jsx
│   ├── pages/            # Page components
│   │   ├── Home.jsx
│   │   ├── Home.css
│   │   ├── Login.jsx
│   │   ├── Signup.jsx
│   │   └── Auth.css
│   ├── services/         # API services
│   │   ├── api.js
│   │   └── authService.js
│   ├── App.jsx          # Main app component with routing
│   ├── App.css          # Global styles
│   └── main.jsx         # Entry point
├── .env                 # Environment variables
└── package.json
```

## Available Scripts

### Development
```bash
npm run dev
```
Starts the development server at `http://localhost:5173`

### Build
```bash
npm run build
```
Builds the app for production to the `dist` folder

### Preview
```bash
npm run preview
```
Preview the production build locally

### Lint
```bash
npm run lint
```
Run ESLint to check code quality

## Features Implemented

### Authentication Flow
1. **Signup Page** (`/signup`)
   - Username, email, and password fields
   - Password validation (min 8 characters)
   - Error handling
   - Link to login page

2. **Login Page** (`/login`)
   - Username and password fields
   - Error handling
   - Link to signup page

3. **Home Page** (`/`)
   - Protected route (requires authentication)
   - Displays welcome message
   - Shows user information
   - Logout button

### Authentication Context
- Global authentication state
- Auto-login on page refresh
- Token storage in localStorage
- Automatic redirect on unauthorized access

### API Integration
- Axios instance with interceptors
- Automatic token injection
- Error handling
- Redirect to login on 401 errors

## Environment Variables

- `VITE_API_BASE_URL` - Backend API URL (default: http://localhost:8080/api)

## Testing the Application

### 1. Start Backend Server
```bash
cd ../backend
./bin/api
```

### 2. Start Frontend Server
```bash
npm run dev
```

### 3. Test Signup
1. Navigate to `http://localhost:5173`
2. You'll be redirected to `/login`
3. Click "Sign up" link
4. Fill in the form with:
   - Username: testuser
   - Email: test@example.com
   - Password: password123
5. Click "Sign Up"
6. You should be redirected to the home page

### 4. Test Login
1. Click "Logout"
2. You'll be redirected to `/login`
3. Enter your credentials
4. Click "Log In"
5. You should be redirected to the home page

## Next Steps

### Planned Features
1. **Post Feed** - Display posts with pagination
2. **Post Submission** - Create new posts (URL/text)
3. **Comments** - Threaded comment system
4. **Voting** - Upvote/downvote functionality
5. **Sorting** - Sort by new, top, best
6. **Search** - Search posts

## Styling

The application uses a Hacker News-inspired color scheme:
- Primary color: `#ff6600` (orange)
- Background: `#f6f6ef` (light beige)
- Clean, minimal design

## License

MIT

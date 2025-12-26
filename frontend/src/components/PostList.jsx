import { useState, useEffect, useRef, useCallback } from 'react';
import Post from './Post';
import { postService } from '../services/postService';
import './PostList.css';

const PostList = ({ sort = 'new', searchQuery = '' }) => {
  const [posts, setPosts] = useState([]);
  const [page, setPage] = useState(1);
  const [loading, setLoading] = useState(false);
  const [hasMore, setHasMore] = useState(true);
  const [error, setError] = useState('');
  const observer = useRef();

  // Load posts
  const loadPosts = useCallback(async (pageNum, reset = false) => {
    if (loading) return;

    setLoading(true);
    setError('');

    try {
      let data;
      if (searchQuery) {
        // Search mode
        data = await postService.searchPosts(searchQuery, pageNum, 20);
      } else {
        // Normal browse mode
        data = await postService.getPosts(pageNum, 20, sort);
      }

      setPosts(prev => reset ? data.posts || [] : [...prev, ...(data.posts || [])]);
      setHasMore((data.posts?.length || 0) === 20);
    } catch (err) {
      setError(err.response?.data?.error || 'Failed to load posts');
    } finally {
      setLoading(false);
    }
  }, [sort, loading, searchQuery]);

  // Reset when sort or search changes
  useEffect(() => {
    setPosts([]);
    setPage(1);
    setHasMore(true);
    loadPosts(1, true);
  }, [sort, searchQuery]);

  // Initial load
  useEffect(() => {
    if (page === 1 && posts.length === 0) {
      loadPosts(1, true);
    }
  }, []);

  // Intersection observer callback
  const lastPostRef = useCallback(node => {
    if (loading) return;
    if (observer.current) observer.current.disconnect();

    observer.current = new IntersectionObserver(entries => {
      if (entries[0].isIntersecting && hasMore) {
        setPage(prevPage => {
          const nextPage = prevPage + 1;
          loadPosts(nextPage);
          return nextPage;
        });
      }
    });

    if (node) observer.current.observe(node);
  }, [loading, hasMore]);

  if (error && posts.length === 0) {
    return <div className="error-message">{error}</div>;
  }

  if (posts.length === 0 && !loading) {
    return (
      <div className="no-posts">
        {searchQuery ? `No posts found for "${searchQuery}"` : 'No posts yet'}
      </div>
    );
  }

  return (
    <div className="post-list">
      {searchQuery && posts.length > 0 && (
        <div className="search-results-header">
          Search results for "{searchQuery}" ({posts.length}+ results)
        </div>
      )}

      {posts.map((post, index) => {
        const isLastPost = index === posts.length - 1;
        return (
          <div key={post.id} ref={isLastPost ? lastPostRef : null}>
            <Post post={post} rank={index + 1} />
          </div>
        );
      })}

      {loading && (
        <div className="loading-indicator">
          <div className="spinner"></div>
          <span>Loading more posts...</span>
        </div>
      )}

      {!hasMore && posts.length > 0 && (
        <div className="end-message">No more posts</div>
      )}
    </div>
  );
};

export default PostList;

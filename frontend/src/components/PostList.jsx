import { useState, useEffect, useRef, useCallback } from 'react';
import Post from './Post';
import { postService } from '../services/postService';
import './PostList.css';

const PostList = ({ sort = 'new', searchQuery = '' }) => {
  const [posts, setPosts] = useState([]);
  const [loading, setLoading] = useState(false);
  const [hasMore, setHasMore] = useState(true);
  const [error, setError] = useState('');
  const observer = useRef();
  const loadingRef = useRef(false);
  const pageRef = useRef(1);

  // Load posts
  const loadPosts = useCallback(async (pageNum, reset = false) => {
    if (loadingRef.current) return;

    loadingRef.current = true;
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
      loadingRef.current = false;
    }
  }, [sort, searchQuery]);

  // Reset when sort or search changes
  useEffect(() => {
    let cancelled = false;

    const loadInitialPosts = async () => {
      // Scroll to top when changing sort or search
      window.scrollTo(0, 0);

      setPosts([]);
      pageRef.current = 1;
      setHasMore(true);
      loadingRef.current = false; // Reset loading flag

      if (!cancelled) {
        await loadPosts(1, true);
      }
    };

    loadInitialPosts();

    return () => {
      cancelled = true;
      loadingRef.current = false; // Reset loading flag on cleanup
    };
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [sort, searchQuery]);

  // Intersection observer callback
  const lastPostRef = useCallback(node => {
    if (loading) return;
    if (observer.current) observer.current.disconnect();

    observer.current = new IntersectionObserver(entries => {
      if (entries[0].isIntersecting && hasMore && !loadingRef.current) {
        pageRef.current += 1;
        loadPosts(pageRef.current);
      }
    }, {
      rootMargin: '200px', // Start loading when within 200px of the last post
      threshold: 0.1
    });

    if (node) observer.current.observe(node);
    // eslint-disable-next-line react-hooks/exhaustive-deps
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

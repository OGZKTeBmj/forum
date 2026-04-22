import './index.css';
import PostWidget from '../../widgets/PostWidget';
import AuthForm from '../../widgets/AuthForm';
import Header from '../../widgets/Header';
import { useEffect, useRef } from 'react';
import { usePostsStore } from '../../store/posts';
import { useUserStore } from '../../store/user';

export default function MainPage() {
  const posts = usePostsStore(state => state.posts);
  const isLoading = usePostsStore(state => state.isLoading);
  const hasMore = usePostsStore(state => state.hasMore);
  const getPostsRequest = usePostsStore(state => state.getPostsRequest);
  const voteRequest = usePostsStore(state => state.voteRequest);
  const deleteVoteRequest = usePostsStore(state => state.deleteVoteRequest)

  const vote = (postId: number, value: number) => {
    if (value === 0) deleteVoteRequest(postId)
    else voteRequest(postId, value)
  }

  const isAuth = useUserStore(state => state.isAuth);

  const bottomRef = useRef<HTMLDivElement | null>(null);

  useEffect(() => {
    getPostsRequest();
  }, [getPostsRequest]);

  useEffect(() => {
    if (!hasMore) return;

    const el = bottomRef.current;
    if (!el) return;

    const observer = new IntersectionObserver(
      (entries) => {
        const entry = entries[0];
        if (entry.isIntersecting && !isLoading) {
          getPostsRequest();
        }
      },
      {
        root: null,
        rootMargin: '200px',
        threshold: 0.1,
      }
    );
    observer.observe(el);

    return () => {
      observer.disconnect();
    };
  }, [posts.length, isLoading, hasMore, getPostsRequest]);

  return (
    <div className="mainPage">
      <Header/>

      <main className="mainContent">
        <div className="contentWrapper">
          <div className="postsSection">
            {posts.length === 0 ? (
              <div className="emptyState">
                <h3>Пока нет постов</h3>
              </div>
            ) : (
              <div className="postsGrid">
                {posts.map((post) => (
                  <PostWidget
                    key={post.id}
                    id={post.id}
                    timeStamp={post.timestamp}
                    title={post.title}
                    content={post.content}
                    author={post.author}
                    commentsCount={post.comments_count}
                    initialRating={post.rating}
                    initialVote={post.user_vote}
                    voteF={(value) => {vote(post.id, value)}}
                  />
                ))}
              </div>
            )}

            <div ref={bottomRef} style={{ height: 1, width: '100%' }} />

            {isLoading && (
              <div className="loadingMore">Загрузка...</div>
            )}
            {!hasMore && posts.length > 0 && (
              <div className="noMorePosts">Больше постов нет</div>
            )}
          </div>

          <aside className="sidebar">
            {!isAuth && <AuthForm />}
          </aside>
        </div>
      </main>
    </div>
  );
}
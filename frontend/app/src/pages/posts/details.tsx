import './index.css';
import PostWidget from '../../widgets/PostWidget';
import AuthForm from '../../widgets/AuthForm';
import Header from '../../widgets/Header';
import { useLocation } from 'react-router-dom';
import { useUserStore } from '../../store/user';
import { usePostsStore } from '../../store/posts';

export default function PostDetailsPage() {
  const location = useLocation();
  const isAuth = useUserStore(state => state.isAuth);
  const post = location.state?.post;

  const voteRequest = usePostsStore(state => state.voteRequest);
  const deleteVoteRequest = usePostsStore(state => state.deleteVoteRequest)

  const vote = (postId: number, value: number) => {
    if (value === 0) deleteVoteRequest(postId)
    else voteRequest(postId, value)
  }
  console.log(post, location.state)
  if (!post) {
    return (
      <div className="mainPage">
        <Header />
        <main className="mainContent">
          <div className="contentWrapper">
            <div className="postsSection">
              <div className="emptyState">
                <h3>Пост не найден</h3>
              </div>
            </div>
            <aside className="sidebar">{!isAuth && <AuthForm />}</aside>
          </div>
        </main>
      </div>
    );
  }

  return (
    <div className="mainPage">
      <Header />

      <main className="mainContent">
        <div className="contentWrapper">
          <div className="postsSection">
            <PostWidget
              id={post.id}
              timeStamp={post.timeStamp}
              title={post.title}
              content={post.content}
              author={post.author}
              commentsCount={post.commentsCount}
              initialRating={post.initialRating}
              initialVote={post.initialVote}
              fullView={post.fullView}
              voteF={(value) => {vote(post.id, value)}}
            />
          </div>

          <aside className="sidebar">{!isAuth && <AuthForm />}</aside>
        </div>
      </main>
    </div>
  );
}
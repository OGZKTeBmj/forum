import { useState } from 'react';
import { useUserStore } from '../../store/user';
import { useNavigate } from 'react-router-dom';
import './index.css';
import MDEditor from '@uiw/react-md-editor';
import CommentsWidget from '../CommentsWidget';
import type { Profile } from '../../models/models';

import defaultAvatar from '../../../assets/icons/default_avatar.png';

export type PostProps = {
  id: number;
  timeStamp: string;
  title: string;
  content?: string | null;
  author: Profile;
  commentsCount: number;
  initialRating?: number;
  initialVote?: number;
  voteF: (value: number) => void;
  fullView?: boolean;
};

function PostWidget(props: PostProps) {
  const {
    id,
    timeStamp,
    title,
    content,
    author,
    commentsCount,
    initialRating = 0,
    initialVote = 0,
    voteF,
    fullView = false,
  } = props;

  const isAuth = useUserStore(state => state.isAuth);
  const navigate = useNavigate();

  const safeContent = content ?? '';

  const [rating, setRating] = useState(initialRating);
  const [vote, setVote] = useState<'none' | 'up' | 'down'>(() => {
    if (initialVote === 1) return 'up';
    if (initialVote === -1) return 'down';
    return 'none';
  });
  const [authMessage, setAuthMessage] = useState<string | null>(null);

  const showAuthMessage = () => {
    setAuthMessage('Войдите или зарегистрируйтесь, чтобы голосовать');
    setTimeout(() => setAuthMessage(null), 2500);
  };

  const createVote = (value: 'none' | 'up' | 'down') => {
    setVote(value);
    if (value === 'none') voteF(0);
    else if (value === 'down') voteF(-1);
    else voteF(1);
  };

  const handleUpvote = () => {
    if (!isAuth) {
      showAuthMessage();
      return;
    }
    if (vote === 'up') {
      setRating(rating - 1);
      createVote('none');
    } else {
      createVote('up');
      if (vote === 'none') setRating(rating + 1);
      else setRating(rating + 2);
    }
  };

  const handleDownvote = () => {
    if (!isAuth) {
      showAuthMessage();
      return;
    }
    if (vote === 'down') {
      setRating(rating + 1);
      createVote('none');
    } else {
      createVote('down');
      if (vote === 'none') setRating(rating - 1);
      else setRating(rating - 2);
    }
  };

  const authorAvatar =
    author.avatar?.thumbnail?.trim()
      ? author.avatar.thumbnail
      : defaultAvatar;

  const handleClickPost = () => {
    if (!fullView) {
      navigate(`/posts/${id}`, {state: {post: {
          id: id,
          timeStamp: timeStamp,
          title: title,
          content: content,
          author: {
            name: author.name,
            avatar: author.avatar
          },
          commentsCount: commentsCount,
          initialRating: initialRating | 0,
          initialVote: initialVote | 0,
          fullView: true,
      }}});
    }
  };

  return (
    <div
      className="postWidget"
      onClick={handleClickPost}
      style={{ cursor: !fullView ? 'pointer' : 'default' }}
    >
      <div className="mainBlock">
        <section className="headerSection">
          <div className="authorBlock">
            <img
              className="authorAvatar"
              src={authorAvatar}
              alt={`Аватар ${author.name}`}
              loading="lazy"
            />
            <span className="authorName">{author.name}</span>
          </div>
          <span className="timeText">{timeStamp}</span>
        </section>

        <section className="titleSection">
          <h2 className="title">{title}</h2>
        </section>

        {fullView && (<section className="contentSection">
          {safeContent ? (
            <MDEditor.Markdown
              className="myMarkDownPreview"
              source={safeContent}
              style={{
                backgroundColor: 'transparent',
                color: 'var(--text-secondary)',
                fontWeight: 400,
                lineHeight: 1.6,
              }}
            />
          ) : (
            <div className="emptyContent">Нет содержимого</div>
          )}
        </section>)}

        <section className="bottomActions">
          <div>
                <button
                  className="toggleButton commentButton"
                >
                  <svg
                    xmlns="http://www.w3.org/2000/svg"
                    width="24"
                    height="24"
                    viewBox="0 0 24 24"
                    fill="rgba(255, 255, 255, 0.36)"
                  >
                    <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z" />
                    <rect x="7" y="7.5" width="11" height="1.6" fill="var(--bg-elevated)" />
                    <rect x="7" y="11.5" width="11" height="1.6" fill="var(--bg-elevated)" />
                  </svg>
                  <span className="commentsCount">{commentsCount}</span>
                </button>
              </div>

              <div className="ratingSection ratingSection--bottom">
                <button
                  className={`ratingButton ratingButton--up${
                    vote === 'up' ? ' ratingButton--voted' : ''
                  }`}
                  onClick={handleUpvote}
                >
                  ▲
                </button>

                <span className="ratingValue">{rating}</span>

                <button
                  className={`ratingButton ratingButton--down${
                    vote === 'down' ? ' ratingButton--voted' : ''
                  }`}
                  onClick={handleDownvote}
                >
                  ▼
                </button>
              </div>
            </section>

        {authMessage && <div className="authMessage">{authMessage}</div>}
        {fullView && <CommentsWidget postId={id} />}
      </div>
    </div>
  );
}

export default PostWidget;
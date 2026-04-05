import { useEffect, useState } from 'react';
import { usePostsStore } from '../../store/posts';
import { useUserStore } from '../../store/user';
import defaultAvatar from '../../../assets/icons/default_avatar.png';
import './index.css';

type CommentsProps = { postId: number };

function CommentsWidget({ postId }: CommentsProps) {
  const commentsStore = usePostsStore(state => state.comments);
  const loadComments = usePostsStore(state => state.getCommentsRequest);
  const createComment = usePostsStore(state => state.createCommentRequest);
  const isAuth = useUserStore(state => state.isAuth);

  const postCommentsObj = commentsStore[postId];
  const comments = postCommentsObj?.batch ?? [];
  const isLoading = postCommentsObj?.isLoading ?? false;
  const hasMore = postCommentsObj?.hasMore ?? true;

  const [text, setText] = useState('');
  const [isSending, setIsSending] = useState(false);

  useEffect(() => {
    if (!postCommentsObj) {
      loadComments(postId);
    }
  }, [postId, postCommentsObj, loadComments]);

  const submitHandler = async () => {
    if (!text.trim()) return;
    setIsSending(true);
    await createComment(postId, text.trim());
    setText('');
    setIsSending(false);
  };

  return (
    <section className="commentsSection">
      <div className="commentsHeader">Комментарии</div>

      {!comments.length && !isLoading && (
        <div className="commentAuthHint">Комментариев пока нет</div>
      )}

      <div className="commentsList">
        {comments.map(comment => (
          <div key={comment.id} className="commentItem">
            <div className="commentContent">
              <div className="commentHeader">
                <img
                  className="authorAvatar"
                  src={ comment.author.avatar?.thumbnail?.trim()
                      ? comment.author.avatar.thumbnail
                      : defaultAvatar}
                  alt={`Аватар ${comment.author.name}`}
                  loading="lazy"
                />
                <span className="commentAuthor">{
                  comment.author.name
                  }</span>
                <span className="commentTime">{comment.timestamp}</span>
              </div>
              <div className="commentText">{comment.content}</div>
            </div>
          </div>
        ))}
      </div>

      {hasMore && (
        <button
          className="loadMoreComments small"
          disabled={isLoading}
          onClick={() => loadComments(postId)}
        >
          {isLoading ? 'Загрузка...' : 'Загрузить ещё'}
        </button>
      )}

      {isAuth ? (
        <>
          <div className="commentsDivider" />
          <textarea
            className="commentInput"
            value={text}
            onChange={e => setText(e.target.value)}
            placeholder="Написать комментарий..."
            disabled={isSending}
          />
          <div className="commentFormActions">
            <button
              className="commentSubmit"
              onClick={submitHandler}
              disabled={!text.trim() || isSending}
            >
              Отправить
            </button>
          </div>
        </>
      ) : (
        <div className="commentAuthHint">
          Только авторизованные пользователи могут оставлять комментарии
        </div>
      )}
    </section>
  );
}

export default CommentsWidget;

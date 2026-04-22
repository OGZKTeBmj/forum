import React, { useState, useEffect, useCallback } from 'react';
import './index.css';
import AuthForm from '../../widgets/AuthForm';
import { useUserStore } from '../../store/user';
import defaultAvatar from '../../../assets/icons/default_avatar.png';
import headerLogo from '../../../assets/icons/header_logo.png';

const ALLOWED_TYPES = [
  'image/png',
  'image/jpeg',
  'image/jpg',
  'image/webp',
];

const MAX_FILE_SIZE = 5 * 1024 * 1024; // 5 MB (nginx client_max_body_size 5m)

const Header: React.FC = () => {
  const isAuth = useUserStore(state => state.isAuth);
  const userName = useUserStore(state => state.name);
  const avatar = useUserStore(state => state.avatar);
  const logout = useUserStore(state => state.logout);
  const uploadAvatar = useUserStore(state => state.uploadAvatar);

  const [isMobile, setIsMobile] = useState(window.innerWidth <= 1024);
  const [showAuthForm, setShowAuthForm] = useState(false);

  const [showAvatarPreview, setShowAvatarPreview] = useState(false);
  const [localPreview, setLocalPreview] = useState<string | null>(null);

  const [isUploading, setIsUploading] = useState(false);
  const [isDragOver, setIsDragOver] = useState(false);
  const [uploadError, setUploadError] = useState<string | null>(null);

  useEffect(() => {
    const handler = () => setIsMobile(window.innerWidth <= 1024);
    window.addEventListener('resize', handler);
    return () => window.removeEventListener('resize', handler);
  }, []);

  const handleLoginClick = () => setShowAuthForm(true);
  const handleCloseAuthForm = () => setShowAuthForm(false);

  const handleAvatarClick = () => {
    setUploadError(null);
    setShowAvatarPreview(true);
  };

  const closeAvatarPreview = () => {
    if (isUploading) return;
    setShowAvatarPreview(false);
    setLocalPreview(null);
    setUploadError(null);
  };

  const processFile = async (file: File) => {
    setUploadError(null);

    if (!ALLOWED_TYPES.includes(file.type)) {
      setUploadError('Поддерживаются только PNG, JPG, JPEG или WEBP');
      return;
    }

    if (file.size > MAX_FILE_SIZE) {
      setUploadError('Файл слишком большой. Максимальный размер — 5 MB');
      return;
    }

    setIsUploading(true);

    const reader = new FileReader();
    reader.onload = () => setLocalPreview(reader.result as string);
    reader.readAsDataURL(file);

    try {
      await uploadAvatar(file);
      window.location.reload();
    } catch (err) {
      console.error('Ошибка загрузки аватара:', err);
      setUploadError('Не удалось загрузить аватар 😢');
      setIsUploading(false);
    }
  };

  const handleAvatarChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (file) processFile(file);
  };

  const handleDrop = useCallback((e: React.DragEvent) => {
    e.preventDefault();
    setIsDragOver(false);
    const file = e.dataTransfer.files[0];
    if (file) processFile(file);
  }, []);

  const smallAvatar = avatar?.thumbnail || defaultAvatar;
  const largeAvatar = localPreview || avatar?.original || defaultAvatar;

  return (
    <header className="mainHeader">
      <div className="headerContent">
        {/* Logo */}
        <img
          src={headerLogo}
          alt="logo"
          className="headerLogo"
          width={176}
          height={53}
        />

        {isAuth ? (
          <div className="profileHeader">
            <span className="profileName">{userName}</span>

            <img
              className="profileAvatar"
              src={smallAvatar}
              onClick={handleAvatarClick}
              alt="avatar"
            />

            <a href="/create-post" className="createPostBtn" aria-label="Создать пост">
              <svg className="plusIcon" viewBox="0 0 24 24" aria-hidden="true">
                <line x1="12" y1="4" x2="12" y2="20" />
                <line x1="4" y1="12" x2="20" y2="12" />
              </svg>
            </a>

            <button className="primaryBtn" onClick={logout}>Выйти</button>
          </div>
        ) : isMobile ? (
          <div className="profileHeader">
            {!showAuthForm && (
              <button className="primaryBtn" onClick={handleLoginClick}>
                Войти
              </button>
            )}
            {showAuthForm && (
              <div className="authModal">
                <div className="authFormWrapper">
                  <AuthForm />
                  <button className="authCloseBtn" onClick={handleCloseAuthForm}>×</button>
                </div>
              </div>
            )}
          </div>
        ) : null}
      </div>

      {showAvatarPreview && (
        <div
          className="avatarModal"
          onDragOver={(e) => {
            e.preventDefault();
            setIsDragOver(true);
          }}
          onDragLeave={() => setIsDragOver(false)}
          onDrop={handleDrop}
        >
          <div className={`avatarModalContent zoomIn ${isDragOver ? 'dragOver' : ''}`}>
            <button className="avatarCloseBtn" onClick={closeAvatarPreview}>×</button>

            <img
              className={`avatarLarge ${isUploading ? 'uploading' : ''}`}
              src={largeAvatar}
              alt="avatar"
            />

            <label className="primaryBtn">
              Загрузить
              <input
                type="file"
                accept="image/png,image/jpeg,image/jpg,image/webp"
                hidden
                onChange={handleAvatarChange}
              />
            </label>

            {isDragOver && <span className="dragHint">Отпускай! 😎</span>}
            {isUploading && <span className="uploadHint">Загрузка...</span>}
            {uploadError && <span className="uploadError">{uploadError}</span>}
          </div>
        </div>
      )}
    </header>
  );
};

export default Header;
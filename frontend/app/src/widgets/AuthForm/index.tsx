import './index.css';
import React, { useState } from 'react';
import type { FormEvent } from 'react';
import { useUserStore } from '../../store/user';
import OpenEye from '../../../assets/icons/open.svg';
import HideEye from '../../../assets/icons/hide.svg';

const AuthForm: React.FC = () => {
  const [shown, setShown] = useState(false);
  const [mode, setMode] = useState<'login' | 'register'>('login');
  const [password, setPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const [name, setName] = useState('');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);
  const { login, registration } = useUserStore();

  const handleSubmit = async (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    setError('');

    let msg: string = "";
    if (name.trim().length < 3) msg = 'Имя должно содержать минимум 3 символа';
    else if (password.length < 8) msg = 'Пароль должен содержать минимум 8 символов';
    else if (mode === 'register' && password !== confirmPassword) msg = 'Пароли не совпадают';

    if (msg) return setError(msg);

    setLoading(true);
    try {
      mode === 'login'
        ? await login(name, password)
        : await registration(name, password);

      setName('');
      setPassword('');
      setConfirmPassword('');
    } catch (err: any) {
      let errorMessage = err?.response?.data?.message || 'Произошла ошибка';
      if (errorMessage === "invalid credentails") {
        errorMessage = "неверные имя или пароль"
      }
      setError(errorMessage);
    } finally {
      setLoading(false);
    }
  };

  const switchMode = () => {
    setMode(prev => prev === 'login' ? 'register' : 'login');
    setError('');
    setName('');
    setPassword('');
    setConfirmPassword('');
  };

  return (
    <div className='wrapper'>
      <div className='authForm'>
        <p className="auth-title">{mode === 'login' ? 'Войти' : 'Регистрация'}</p>

        <form onSubmit={handleSubmit}>
          
          <section>
            <input
              placeholder='Имя'
              type="text"
              value={name}
              onChange={e => setName(e.target.value)}
              required
              disabled={loading}
            />
          </section>

          <section className="password-wrapper">
            <input
              placeholder='Пароль'
              type={shown ? "text" : "password"}
              value={password}
              onChange={e => setPassword(e.target.value)}
              required
              disabled={loading}
            />
            <img
              src={shown ? HideEye : OpenEye}
              alt="toggle visibility"
              className="eye-icon"
              onClick={() => !loading && setShown(!shown)}
            />
          </section>

          {mode === 'register' && (
            <section className="password-wrapper">
              <input
                placeholder='Повторите пароль'
                type={shown ? "text" : "password"}
                value={confirmPassword}
                onChange={e => setConfirmPassword(e.target.value)}
                required
                disabled={loading}
              />
              <img
                src={shown ? HideEye : OpenEye}
                alt="toggle visibility"
                className="eye-icon"
                onClick={() => !loading && setShown(!shown)}
              />
            </section>
          )}

          {error && <div className="error-message">{error}</div>}

          <button type="submit" className='b1' disabled={loading}>
            {loading ? 'Загрузка...' : (mode === 'login' ? 'Войти' : 'Зарегистрироваться')}
          </button>
        </form>

        <button onClick={switchMode} className='b2' disabled={loading}>
          {mode === 'login' ? 'Создать аккаунт' : 'У меня есть аккаунт'}
        </button>
      </div>
    </div>
  );
};

export default AuthForm;

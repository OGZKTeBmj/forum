import { Routes, Route } from 'react-router-dom';
import { useEffect } from 'react';
import MainPage from './pages/posts/main';
import CreatePostPage from './pages/createPost';
import PostDetailPage from './pages/posts/details';

import { useUserStore } from './store/user';

function App() {
  const checkAuth = useUserStore((state) => state.checkAuth)

  useEffect(() => {
    checkAuth();
  }, [checkAuth]);

  return (
    <Routes>
      <Route path="/" element={<MainPage />} />
      <Route path="/create-post" element={<CreatePostPage />} />
      <Route path="/posts/:id" element={<PostDetailPage />} />
    </Routes>
  );
}

export default App;

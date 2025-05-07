import { createBrowserRouter, RouterProvider } from 'react-router-dom';
import { Provider } from 'react-redux';
import { store } from './store/store';
import Dashboard from './pages/Dashboard';
import Transactions from './pages/Transactions';
import Categories from './pages/Categories';
import Reports from './pages/Reports';
import Login from './pages/Login';
import Register from './pages/Register';
import PrivateRoute from './components/PrivateRoute';

const router = createBrowserRouter([
  {
    path: '/login',
    element: <Login />
  },
  {
    path: '/register',
    element: <Register />
  },
  {
    path: '/',
    element: <PrivateRoute><Dashboard /></PrivateRoute>
  },
  {
    path: '/transactions',
    element: <PrivateRoute><Transactions /></PrivateRoute>
  },
  {
    path: '/categories',
    element: <PrivateRoute><Categories /></PrivateRoute>
  },
  {
    path: '/reports',
    element: <PrivateRoute><Reports /></PrivateRoute>
  }
], {
  future: {
    v7_startTransition: true
  }
});

const App = () => {
  return (
    <Provider store={store}>
      <RouterProvider router={router} />
    </Provider>
  );
};

export default App; 
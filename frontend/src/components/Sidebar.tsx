import React from 'react';
import { Link, useLocation } from 'react-router-dom';
import logo from '../assets/logo.svg';

const Sidebar: React.FC = () => {
  const location = useLocation();

  const navItems = [
    { name: 'Dashboard', path: '/' },
    { name: 'Launchpad', path: '/launchpad' },
    { name: 'Settings', path: '/settings' },
  ];

  return (
    <aside className="w-64 bg-gray-900 shadow-lg flex flex-col">
      <div className="p-6 flex items-center justify-center">
        <img src={logo} alt="LMNinja Logo" className="h-12 w-12 mr-3" />
        <h1 className="text-2xl font-semibold text-white">LMNinja</h1>
      </div>
      <nav className="flex-grow mt-8">
        <ul>
          {navItems.map((item) => (
            <li key={item.name}>
              <Link
                to={item.path}
                className={`block py-3 px-6 text-lg transition-colors duration-200 ${location.pathname === item.path ? 'bg-gray-700 text-white border-l-4 border-purple-500' : 'text-gray-300 hover:bg-gray-700'}`}
              >
                {item.name}
              </Link>
            </li>
          ))}
        </ul>
      </nav>
      <div className="p-6 text-gray-500 text-sm">
        <p>&copy; 2023 LMNinja</p>
      </div>
    </aside>
  );
};

export default Sidebar;
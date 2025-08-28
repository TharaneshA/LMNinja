import React from 'react';

const Dashboard: React.FC = () => {
  return (
    <div className="dashboard p-6 bg-gray-800 rounded-lg shadow-xl">
      <h1 className="text-4xl font-bold text-white mb-6">Welcome to LMNinja!</h1>
      <p className="text-lg text-gray-300 mb-8">
        Your ultimate tool for red teaming and evaluating Large Language Models.
      </p>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        <div className="stat-card bg-gray-700 p-5 rounded-lg shadow-md">
          <h2 className="text-xl font-semibold text-purple-400 mb-2">Security Score</h2>
          <p className="text-3xl font-bold text-green-400">85/100</p>
          <p className="text-gray-400">Based on last 7 days scans</p>
        </div>

        <div className="stat-card bg-gray-700 p-5 rounded-lg shadow-md">
          <h2 className="text-xl font-semibold text-purple-400 mb-2">Recent Scans</h2>
          <p className="text-3xl font-bold text-blue-400">12</p>
          <p className="text-gray-400">Total scans performed</p>
        </div>

        <div className="stat-card bg-gray-700 p-5 rounded-lg shadow-md">
          <h2 className="text-xl font-semibold text-purple-400 mb-2">Vulnerabilities Found</h2>
          <p className="text-3xl font-bold text-red-400">3</p>
          <p className="text-gray-400">Critical issues detected</p>
        </div>
      </div>

      <div className="mt-10">
        <h2 className="text-2xl font-bold text-white mb-4">Quick Actions</h2>
        <div className="flex space-x-4">
          <button className="bg-purple-600 hover:bg-purple-700 text-white font-bold py-3 px-6 rounded-lg transition duration-300">
            Start New Scan
          </button>
          <button className="bg-blue-600 hover:bg-blue-700 text-white font-bold py-3 px-6 rounded-lg transition duration-300">
            View Reports
          </button>
        </div>
      </div>
    </div>
  );
};

export default Dashboard;
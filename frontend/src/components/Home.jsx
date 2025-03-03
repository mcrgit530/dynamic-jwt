import { Link } from "react-router-dom";

function Home() {
  return (
    <div className="h-screen flex flex-col items-center justify-center bg-gray-100">
      <h1 className="text-4xl font-bold mb-6 text-blue-600">Welcome</h1>
      <div className="space-x-4">
        <Link to="/signup">
          <button className="px-6 py-2 bg-green-500 text-white font-semibold rounded-lg shadow-md hover:bg-green-600">
            Sign Up
          </button>
        </Link>
        <Link to="/signin">
          <button className="px-6 py-2 bg-blue-500 text-white font-semibold rounded-lg shadow-md hover:bg-blue-600">
            Sign In
          </button>
        </Link>
      </div>
    </div>
  );
}

export default Home;

// src/app/page.tsx
import Link from 'next/link';

export default function Home() {
  // The requested dummy profile UUID
  const demoUserId = "9db2d431-cee7-43aa-b237-6ca80932bfae";

  return (
    <main className="min-h-screen bg-slate-50 flex flex-col items-center justify-center p-6 font-sans">
      <div className="text-center max-w-lg">
        <h1 className="text-4xl font-extrabold text-slate-900 tracking-tight mb-4">
          Contest Rating Engine
        </h1>
        <p className="text-lg text-slate-600 mb-8">
          The backend engine is live. The frontend is fully connected to the Go REST API.
        </p>
        
        <Link 
          href="/profile?id=9db2d431-cee7-43aa-b237-6ca80932bfae"
          className="inline-flex items-center justify-center bg-blue-600 text-white px-8 py-4 rounded-xl font-bold text-lg shadow-lg hover:bg-blue-700 hover:-translate-y-1 transition-all duration-200"
        >
          View Demo Profile
        </Link>
      </div>
    </main>
  );
}
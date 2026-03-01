"use client";

import Link from 'next/link';
import { FileQuestion, Home, ArrowLeft } from 'lucide-react';

export default function NotFound() {
  return (
    <div className="min-h-screen bg-slate-50 flex items-center justify-center px-4 font-sans">
      <div className="max-w-md w-full text-center space-y-8">
        {/* Visual Cue */}
        <div className="flex justify-center">
          <div className="bg-amber-100 p-6 rounded-full">
            <FileQuestion className="w-16 h-16 text-amber-600" />
          </div>
        </div>

        {/* Text Content */}
        <div className="space-y-3">
          <h1 className="text-5xl font-black text-slate-900 tracking-tight">404</h1>
          <h2 className="text-2xl font-bold text-slate-800">Profile Not Found</h2>
          <p className="text-slate-500 leading-relaxed">
            The Rating Engine couldn not t find a user with this ID. 
            It might have been lost in the database or the URL is malformed.
          </p>
        </div>

        {/* Action Buttons */}
        <div className="flex flex-col sm:flex-row gap-4 justify-center pt-4">
          <button 
            onClick={() => window.history.back()}
            className="flex items-center justify-center gap-2 px-6 py-3 bg-white border border-slate-200 text-slate-700 font-bold rounded-xl hover:bg-slate-50 transition-all shadow-sm"
          >
            <ArrowLeft className="w-4 h-4" />
            Go Back
          </button>
          
          <Link 
            href="/"
            className="flex items-center justify-center gap-2 px-6 py-3 bg-blue-600 text-white font-bold rounded-xl hover:bg-blue-700 transition-all shadow-md shadow-blue-200"
          >
            <Home className="w-4 h-4" />
            Return Home
          </Link>
        </div>

        {/* System Footer */}
        <div className="pt-8 border-t border-slate-200">
          <p className="text-xs text-slate-400 uppercase tracking-widest font-bold">
            RankShift Rating Engine v1.0
          </p>
        </div>
      </div>
    </div>
  );
}


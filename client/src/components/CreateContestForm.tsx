// src/components/CreateContestForm.tsx
'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';

const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';

export default function CreateContestForm() {
  const router = useRouter();
  const [name, setName] = useState('');
  const [participants, setParticipants] = useState<number | ''>('');
  const [isSubmitting, setIsSubmitting] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!name || participants === '') return;
    
    setIsSubmitting(true);

    try {
      const res = await fetch(`${API_URL}/api/contests`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          name: name,
          total_participants: Number(participants)
        }),
      });

      if (!res.ok) throw new Error("Failed to create contest.");
      
      router.refresh(); // Instantly refreshes the page to update the dropdown list!
      setName(''); 
      setParticipants('');
    }  catch (error: unknown) {
        if (error instanceof Error) {
            alert(error.message);
        } else {
            alert("Something went wrong.");
        }
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <form onSubmit={handleSubmit} className="bg-white rounded-2xl border border-slate-200 shadow-sm p-6 sm:p-8">
      <h2 className="text-lg font-bold text-slate-900 mb-4">Host a New Contest</h2>
      <div className="flex flex-col md:flex-row gap-4">
        
        <div className="flex-1">
          <label className="block text-xs font-bold text-slate-500 uppercase tracking-wider mb-1">Contest Name</label>
          <input 
            type="text" 
            required
            placeholder="e.g., Visual Vortex 3.0"
            value={name}
            onChange={(e) => setName(e.target.value)}
            className="w-full px-4 py-2 border border-slate-300 rounded-lg focus:ring-2 focus:ring-blue-500 outline-none text-sm text-slate-900"
          />
        </div>

        <div className="w-40">
          <label className="block text-xs font-bold text-slate-500 uppercase tracking-wider mb-1">Total Players</label>
          <input 
            type="number" 
            min="1" 
            required 
            placeholder="e.g., 500"
            value={participants} 
            onChange={(e) => setParticipants(e.target.value === '' ? '' : Number(e.target.value))}
            className="w-full px-4 py-2 border border-slate-300 rounded-lg focus:ring-2 focus:ring-blue-500 outline-none text-sm text-slate-900" 
          />
        </div>

        <div className="flex items-end">
          <button 
            type="submit" 
            disabled={isSubmitting || !name || participants === ''}
            className="h-[38px] px-6 bg-blue-600 text-white font-bold rounded-lg hover:bg-blue-700 disabled:opacity-50 text-sm transition-colors"
          >
            {isSubmitting ? 'Creating...' : 'Create'}
          </button>
        </div>
      </div>
    </form>
  );
}
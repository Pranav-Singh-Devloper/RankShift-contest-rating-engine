'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import { Contest } from '@/lib/types';

const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';

interface AddContestFormProps {
  userId: string;
  availableContests: Contest[];
}

export default function AddContestForm({ userId, availableContests }: AddContestFormProps) {
  const router = useRouter();
  const [contestId, setContestId] = useState('');
  const [rank, setRank] = useState<number | ''>('');
  const [isSubmitting, setIsSubmitting] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!contestId || rank === '') return;
    
    setIsSubmitting(true);

    try {
      const payload = {
        contest_id: contestId,
        results: [{ user_id: userId, rank: Number(rank) }]
      };

      const res = await fetch(`${API_URL}:8080/api/contests/end`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(payload),
      });

      if (!res.ok) throw new Error("Failed to process contest results.");
      
      router.refresh(); 
      setContestId(''); 
      setRank('');
    } catch (error: unknown) {
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
    <form onSubmit={handleSubmit} className="bg-white rounded-2xl border border-slate-200 shadow-sm p-6 sm:p-8 mt-6">
      <h2 className="text-lg font-bold text-slate-900 mb-4">Simulate New Result</h2>
      <div className="flex flex-col md:flex-row gap-4">
        
        {/* The New Dropdown Menu */}
        <div className="flex-1">
          <label className="block text-xs font-bold text-slate-500 uppercase tracking-wider mb-1">Select Contest</label>
          <select 
            required
            value={contestId}
            onChange={(e) => setContestId(e.target.value)}
            className="w-full px-4 py-2 border border-slate-300 rounded-lg focus:ring-2 focus:ring-blue-500 outline-none text-sm text-slate-900 bg-white"
          >
            <option value="" disabled>-- Choose a Contest --</option>
            {availableContests.map((contest) => (
              <option key={contest.id} value={contest.id}>
                {contest.name} ({contest.total_participants} Participants)
              </option>
            ))}
          </select>
        </div>

        <div className="w-32">
          <label className="block text-xs font-bold text-slate-500 uppercase tracking-wider mb-1">Rank</label>
          <input 
            type="number" 
            min="1" 
            required 
            placeholder="e.g. 5"
            value={rank} 
            onChange={(e) => setRank(e.target.value === '' ? '' : Number(e.target.value))}
            className="w-full px-4 py-2 border border-slate-300 rounded-lg focus:ring-2 focus:ring-blue-500 outline-none text-sm text-slate-900" 
          />
        </div>

        <div className="flex items-end">
          <button 
            type="submit" 
            disabled={isSubmitting || !contestId || rank === ''}
            className="h-[38px] px-6 bg-slate-900 text-white font-bold rounded-lg hover:bg-slate-800 disabled:opacity-50 text-sm transition-colors"
          >
            {isSubmitting ? 'Calculating...' : 'Submit'}
          </button>
        </div>
      </div>
    </form>
  );
}
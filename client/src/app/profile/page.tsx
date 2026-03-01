// src/app/profile/page.tsx
import { Contest, UserProfile } from '@/lib/types';
import { getTierColor } from '@/lib/utils';
import { Trophy, TrendingUp, Award, Activity } from 'lucide-react';
import { notFound } from 'next/navigation';
import InteractiveHistory from '@/components/InteractiveHistory';
import AddContestForm from '@/components/AddContestForm';
import CreateContestForm from '@/components/CreateContestForm';

export const dynamic = 'force-dynamic';

const API_URL = 'https://rankshift-contest-rating-engine-1.onrender.com';

async function getProfileData(userId: string): Promise<UserProfile | null> {
  try {
    const res = await fetch(`${API_URL}/api/users/${userId}`, {
      cache: 'no-store',
    });
    if (!res.ok) return null;
    return await res.json();
  } catch (error) {
    console.error("Fetch error:", error);
    return null;
  }
}

async function getAllContests(): Promise<Contest[]> {
  try {
    const res = await fetch(`${API_URL}/api/contests`, { cache: 'no-store' });
    if (!res.ok) return [];
    return await res.json();
  } catch (error) {
    return [];
  }
}

// In Next.js 15, searchParams is also a Promise
type Props = {
  searchParams: Promise<{ id?: string }>;
};

export default async function ProfilePage({ searchParams }: Props) {
  const { id } = await searchParams;

  // If no ID is provided in the URL, show 404 or redirect
  if (!id) {
    notFound();
  }

  const user = await getProfileData(id);
  const availableContests = await getAllContests();

  if (!user) {
    notFound();
  }

  const tierColorClass = getTierColor(user.tier);

  return (
    <main className="min-h-screen bg-slate-50 py-12 px-4 sm:px-6 lg:px-8">
      <div className="max-w-4xl mx-auto space-y-6">
        {/* Profile Header Card */}
        <div className="bg-white rounded-2xl border border-slate-200 p-8 shadow-sm">
          <h1 className="text-3xl font-bold text-slate-900">{user.name}</h1>
          <div className="flex items-center gap-2 mt-2">
            <Trophy className={`w-5 h-5 ${tierColorClass}`} />
            <span className={`font-bold uppercase tracking-wider ${tierColorClass}`}>
              {user.tier}
            </span>
            <span className="text-slate-400 text-sm ml-2">
              • {user.contests_played} Contests
            </span>
          </div>

          <div className="grid grid-cols-2 gap-4 mt-6">
            <div className="bg-blue-50 border border-blue-100 p-4 rounded-xl">
              <p className="text-xs text-blue-600 font-bold uppercase">Current Rating</p>
              <p className="text-3xl font-black text-blue-950">{user.current_rating}</p>
            </div>
            <div className="bg-amber-50 border border-amber-100 p-4 rounded-xl">
              <p className="text-xs text-amber-600 font-bold uppercase">Peak Rating</p>
              <p className="text-3xl font-black text-amber-950">{user.max_rating}</p>
            </div>
          </div>
        </div>

        {/* Chart Section */}
        <div className="bg-white rounded-2xl border border-slate-200 p-8 shadow-sm">
          <h2 className="text-lg font-bold text-slate-900 flex items-center gap-2 mb-4">
            <TrendingUp className="w-5 h-5 text-slate-500" />
            Rating Progression
          </h2>
          <InteractiveHistory data={user.ratingHistory} />
        </div>

        {/* Forms Section */}
        <div className="grid grid-cols-1 gap-6">
          <CreateContestForm />
          <AddContestForm userId={user.id} availableContests={availableContests} />
        </div>
      </div>
    </main>
  );
}
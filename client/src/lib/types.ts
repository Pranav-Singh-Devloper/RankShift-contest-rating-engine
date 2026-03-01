// src/lib/types.ts

export interface Contest {
  id: string;
  name: string;
  date: string;
  total_participants?: number; // Optional based on your Prisma include
}

export interface RatingHistory {
  id: string;
  user_id: string;
  contest_id: string;
  old_rating: number;
  new_rating: number;
  performance_rating: number;
  rank: number;
  percentile: number;
  rating_change: number;
  contest: Contest;
}

export interface UserProfile {
  id: string;
  name: string;
  current_rating: number;
  max_rating: number;
  contests_played: number;
  tier: string;
  ratingHistory: RatingHistory[]; // Lowercase 'r' to match JSON
}
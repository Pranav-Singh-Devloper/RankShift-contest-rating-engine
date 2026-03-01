// src/components/InteractiveHistory.tsx
'use client';

import { useState } from 'react';
import {
  LineChart,
  Line,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
} from 'recharts';
import { RatingHistory } from '@/lib/types';
import { Target, Hash, Percent, TrendingUp } from 'lucide-react';

type ChartPoint = {
  displayName: string;
  rating: number;
  originalData: RatingHistory;
};

export default function InteractiveHistory({
  data,
}: {
  data: RatingHistory[];
}) {
  const [selectedContest, setSelectedContest] =
    useState<RatingHistory | null>(null);

  const safeData = data || [];

  if (safeData.length === 0) {
    return (
      <div className="h-72 flex items-center justify-center text-slate-400">
        No history available.
      </div>
    );
  }

  // Sort + map safely
  const chartData: ChartPoint[] = [...safeData]
    .sort(
      (a, b) =>
        new Date(a.contest.date).getTime() -
        new Date(b.contest.date).getTime()
    )
    .map((entry, index) => ({
      displayName: `Contest ${index + 1}`,
      rating: entry.new_rating,
      originalData: entry,
    }));

  return (
    <div className="mt-6">
      {/* Chart */}
      <div className="h-[300px] w-full cursor-pointer">
        <ResponsiveContainer width="100%" height="100%">
          <LineChart
            data={chartData}
            margin={{ top: 10, right: 10, bottom: 10, left: -20 }}
          >
            <CartesianGrid
              strokeDasharray="3 3"
              vertical={false}
              stroke="#f1f5f9"
            />

            <XAxis
              dataKey="displayName"
              axisLine={false}
              tickLine={false}
              tick={{ fontSize: 12, fill: '#64748b' }}
              dy={10}
            />

            <YAxis
              domain={['auto', 'auto']}
              axisLine={false}
              tickLine={false}
              tick={{ fontSize: 12, fill: '#64748b' }}
            />

            <Tooltip
              cursor={{ stroke: '#cbd5e1', strokeDasharray: '4 4' }}
              content={({ active, payload }) => {
                if (active && payload && payload.length) {
                  const point = payload[0].payload as ChartPoint;
                  return (
                    <div className="bg-slate-900 text-white text-xs px-3 py-2 rounded-md shadow-lg">
                      <p className="font-bold">{point.displayName}</p>
                      <p>Click dot for details</p>
                    </div>
                  );
                }
                return null;
              }}
            />

            <Line
              type="monotone"
              dataKey="rating"
              stroke="#0ea5e9"
              strokeWidth={3}
              dot={{ r: 5, fill: '#0ea5e9', strokeWidth: 2, stroke: '#ffffff' }}
              activeDot={(props) => {
                const point = props.payload as ChartPoint;

                return (
                  <circle
                    cx={props.cx}
                    cy={props.cy}
                    r={8}
                    stroke="#bae6fd"
                    strokeWidth={4}
                    fill="#0ea5e9"
                    style={{ cursor: 'pointer' }}
                    onClick={() =>
                      setSelectedContest(point.originalData)
                    }
                  />
                );
              }}
            />
          </LineChart>
        </ResponsiveContainer>
      </div>

      {/* Details Panel */}
      {selectedContest && (
        <div className="mt-6 p-5 bg-slate-50 border border-slate-200 rounded-xl animate-in fade-in slide-in-from-top-4 duration-300">
          <h3 className="text-sm font-bold text-slate-500 uppercase tracking-wider mb-4 border-b pb-2">
            Contest Mathematics (ID:{' '}
            {selectedContest.contest_id.split('-')[0]}...)
          </h3>

          <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
            <div>
              <p className="text-xs text-slate-500 flex items-center gap-1">
                <Hash className="w-3 h-3" /> Rank
              </p>
              <p className="font-bold text-slate-900">
                {selectedContest.rank}
              </p>
            </div>

            <div>
              <p className="text-xs text-slate-500 flex items-center gap-1">
                <Percent className="w-3 h-3" /> Percentile
              </p>
              <p className="font-bold text-slate-900">
                {(selectedContest.percentile * 100).toFixed(1)}%
              </p>
            </div>

            <div>
              <p className="text-xs text-slate-500 flex items-center gap-1">
                <Target className="w-3 h-3" /> Perf. Rating
              </p>
              <p className="font-bold text-slate-900">
                {selectedContest.performance_rating}
              </p>
            </div>

            <div>
              <p className="text-xs text-slate-500 flex items-center gap-1">
                <TrendingUp className="w-3 h-3" /> Rating Shift
              </p>
              <p className="font-bold text-slate-900">
                {selectedContest.old_rating} →{' '}
                {selectedContest.new_rating}
                <span
                  className={`ml-2 ${
                    selectedContest.rating_change > 0
                      ? 'text-green-600'
                      : 'text-red-600'
                  }`}
                >
                  (
                  {selectedContest.rating_change > 0 ? '+' : ''}
                  {selectedContest.rating_change})
                </span>
              </p>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
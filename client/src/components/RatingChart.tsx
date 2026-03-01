'use client';

import {
  LineChart,
  Line,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
  TooltipProps,
} from 'recharts';
import { RatingHistory } from '@/lib/types';

interface RatingChartProps {
  data: RatingHistory[];
}

/* ---------- Chart Point Type ---------- */
type ChartPoint = {
  name: string;
  rating: number;
  rank: number;
  change: number;
  date: string;
};

/* ===================================================== */
/*   CUSTOM TOOLTIP COMPONENT                           */
/* ===================================================== */


function CustomTooltip({
  active,
  payload,
}: TooltipProps<number, string> & { payload?: { payload: ChartPoint }[] }) {
  if (!active || !payload || payload.length === 0) {
    return null;
  }

  // Safely extract the data from the payload
  const tooltipData = payload[0]?.payload;

  if (!tooltipData) return null;

  const changeColor =
    tooltipData.change > 0
      ? 'text-green-600'
      : tooltipData.change < 0
      ? 'text-red-600'
      : 'text-slate-600';

  const changePrefix = tooltipData.change > 0 ? '+' : '';

  return (
    <div className="bg-white border border-slate-200 p-4 rounded-lg shadow-xl text-sm min-w-[160px]">
      <p className="font-bold text-slate-900 mb-2 border-b pb-1">
        {tooltipData.name}
      </p>

      <div className="flex justify-between items-center mb-1">
        <span className="text-slate-500">Rank:</span>
        <span className="font-semibold text-slate-800">
          {tooltipData.rank}
        </span>
      </div>

      <div className="flex justify-between items-center">
        <span className="text-slate-500">Rating Change:</span>
        <span className={`font-bold ${changeColor}`}>
          {changePrefix}
          {tooltipData.change}
        </span>
      </div>
    </div>
  );
}

/* ===================================================== */
/*   MAIN COMPONENT                                     */
/* ===================================================== */

export default function RatingChart({ data }: RatingChartProps) {
  const safeData = data ?? [];

  if (safeData.length === 0) {
    return (
      <div className="h-72 w-full flex items-center justify-center text-slate-400 border-2 border-dashed border-slate-200 rounded-xl mt-6 bg-slate-50/50">
        No contest history available yet.
      </div>
    );
  }

  /* ---------- Sort + Format Data Safely ---------- */
  const chartData: ChartPoint[] = [...safeData]
    .sort(
      (a, b) =>
        new Date(a.contest.date).getTime() -
        new Date(b.contest.date).getTime()
    )
    .map((entry) => ({
      name: entry.contest.name,
      rating: entry.new_rating,
      rank: entry.rank,
      change: entry.rating_change,
      date: new Date(entry.contest.date).toLocaleDateString(),
    }));

  return (
    <div className="h-[350px] w-full mt-6">
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
            dataKey="name"
            axisLine={false}
            tickLine={false}
            tick={{
              fontSize: 12,
              fill: '#64748b',
              fontWeight: 500,
            }}
            dy={15}
          />

          <YAxis
            domain={['auto', 'auto']}
            axisLine={false}
            tickLine={false}
            tick={{
              fontSize: 12,
              fill: '#64748b',
              fontWeight: 500,
            }}
          />

          <Tooltip
            content={<CustomTooltip />}
            cursor={{
              stroke: '#cbd5e1',
              strokeWidth: 1,
              strokeDasharray: '4 4',
            }}
          />

          <Line
            type="monotone"
            dataKey="rating"
            stroke="#0ea5e9"
            strokeWidth={3}
            dot={{
              r: 4,
              fill: '#0ea5e9',
              strokeWidth: 2,
              stroke: '#ffffff',
            }}
            activeDot={{
              r: 6,
              stroke: '#bae6fd',
              strokeWidth: 4,
            }}
            animationDuration={1200}
          />
        </LineChart>
      </ResponsiveContainer>
    </div>
  );
}
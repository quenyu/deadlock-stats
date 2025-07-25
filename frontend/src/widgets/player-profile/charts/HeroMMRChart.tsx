import { useState } from 'react';
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from "@/shared/ui/card";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/shared/ui/select";
import { LineChart, Line, XAxis, YAxis, Tooltip, ResponsiveContainer, Legend } from 'recharts';

interface MMRPoint {
  match_id: number;
  start_time: number;
  player_score: number;
  rank: number;
}

interface HeroMMRHistory {
  hero_id: number;
  hero_name: string;
  history: MMRPoint[];
}

interface HeroMMRChartProps {
  heroMMRHistory: HeroMMRHistory[];
}

const CustomTooltip = ({ active, payload }: any) => {
  if (active && payload && payload.length) {
    const data = payload[0].payload;
    const date = new Date(data.start_time * 1000).toLocaleDateString();
    return (
      <div className="rounded-lg border bg-background p-2 shadow-sm">
        <p className="text-sm font-bold">{`Rank: ${data.rank}`}</p>
        <p className="text-xs text-muted-foreground">{`Score: ${data.player_score.toFixed(0)}`}</p>
        <p className="text-xs text-muted-foreground">{date}</p>
      </div>
    );
  }
  return null;
};

export const HeroMMRChart = ({ heroMMRHistory }: HeroMMRChartProps) => {
  const [selectedHeroId, setSelectedHeroId] = useState<string | null>(
    heroMMRHistory?.[0]?.hero_id.toString() || null
  );

  if (!heroMMRHistory || heroMMRHistory.length === 0) {
    return (
      <Card>
        <CardHeader>
          <CardTitle>Hero MMR Progress</CardTitle>
          <CardDescription>No data available to show MMR progress for heroes.</CardDescription>
        </CardHeader>
      </Card>
    );
  }

  const selectedHeroData = heroMMRHistory.find(h => h.hero_id.toString() === selectedHeroId);

  console.log(heroMMRHistory);
  return (
    <Card>
      <CardHeader>
        <div className="flex justify-between items-center">
          <div>
            <CardTitle>Hero MMR Progress</CardTitle>
            <CardDescription>Track your MMR changes on specific heroes.</CardDescription>
          </div>
          <Select value={selectedHeroId ?? ''} onValueChange={setSelectedHeroId}>
            <SelectTrigger className="w-[180px]">
              <SelectValue placeholder="Select a hero" />
            </SelectTrigger>
            <SelectContent>
              {heroMMRHistory.map(hero => (
                <SelectItem key={hero.hero_id} value={hero.hero_id.toString()}>
                  {hero.hero_name}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
        </div>
      </CardHeader>
      <CardContent>
        <ResponsiveContainer width="100%" height={300}>
          <LineChart data={selectedHeroData?.history}>
            <XAxis 
              dataKey="start_time" 
              tickFormatter={(time) => new Date(time * 1000).toLocaleDateString()} 
              tick={{ fontSize: 12 }}
            />
            <YAxis 
              domain={['dataMin - 10', 'dataMax + 10']} 
              tick={{ fontSize: 12 }} 
              tickFormatter={(value) => Math.round(value).toString()}
            />
            <Tooltip content={<CustomTooltip />} />
            <Legend />
            <Line type="monotone" dataKey="player_score" name="MMR Score" stroke="#8884d8" strokeWidth={2} dot={false} />
          </LineChart>
        </ResponsiveContainer>
      </CardContent>
    </Card>
  );
}; 
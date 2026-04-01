import { Card, CardContent, CardHeader, CardTitle, CardDescription } from "@/shared/ui/card";
import { Avatar, AvatarImage, AvatarFallback } from "@/shared/ui/avatar";

interface MateStat {
  steam_id: string;
  nickname: string;
  avatar_url: string;
  games: number;
  wins: number;
  win_rate: number;
}

interface BestMatesProps {
  mates: MateStat[];
}

export const BestMates = ({ mates }: BestMatesProps) => {
  const safeToFixed = (value: number | undefined | null, decimals: number = 1): string => {
    if (value === undefined || value === null || isNaN(value)) {
      return '0.0';
    }
    return value.toFixed(decimals);
  };

  if (!mates || mates.length === 0) {
    return null;
  }

  const topMates = mates.slice(0, 5);

  return (
    <Card>
      <CardHeader>
        <CardTitle>Best Mates</CardTitle>
        <CardDescription>Players you have the highest win rate with.</CardDescription>
      </CardHeader>
      <CardContent className="space-y-4">
        {topMates.map((mate) => (
          <div key={mate.steam_id} className="flex items-center justify-between p-2 rounded-lg transition-colors hover:bg-muted/50">
            <div className="flex items-center gap-4">
              <Avatar className="h-10 w-10">
                <AvatarImage src={mate.avatar_url} alt={mate.nickname} />
                <AvatarFallback>{mate.nickname ? mate.nickname.charAt(0) : "?"}</AvatarFallback>
              </Avatar>
              <div>
                <p className="font-semibold text-sm">{mate.nickname}</p>
                <p className="text-xs text-muted-foreground">{mate.games} games played</p>
              </div>
            </div>
            <div className="text-right">
              <p className="font-bold text-green-500 text-sm">{safeToFixed(mate.win_rate, 1)}%</p>
              <p className="text-xs text-muted-foreground">Win Rate</p>
            </div>
          </div>
        ))}
      </CardContent>
    </Card>
  );
}; 
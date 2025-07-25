import { ArrowUpIcon, ArrowDownIcon } from "@radix-ui/react-icons";
import { Card, CardContent, CardHeader } from "@/shared/ui/card";

interface PerformanceSnapshotProps {
  stats: {
    win_rate: number;
    kd_ratio: number;
    total_matches: number;
    performance_dynamics: {
      win_loss: {
        trend: string;
        value: string;
        sparkline: number[];
      };
      kda: {
        trend: string;
        value: string;
        sparkline: number[];
      };
      rank: {
        trend: string;
        value: string;
        sparkline: number[];
      };
    };
    avg_kills_per_match: number;
    avg_deaths_per_match: number;
    avg_assists_per_match: number;
    avg_match_duration: number;
  };
}

export const PerformanceSnapshot = ({ stats }: PerformanceSnapshotProps) => {
  const renderTrend = (trend: string) => {
    if (trend === "up") {
      return <ArrowUpIcon className="text-green-500" />;
    } else if (trend === "down") {
      return <ArrowDownIcon className="text-red-500" />;
    }
    return null;
  };

  return (
    <Card>
      <CardHeader>
        <h3 className="text-lg font-semibold">Performance Snapshot</h3>
      </CardHeader>
      <CardContent className="space-y-6">
        <div className="space-y-4">
          <div className="flex justify-between">
            <div className="space-y-0.5">
              <div className="text-sm text-muted-foreground">Win Rate</div>
              <div className="text-2xl font-bold">{stats.win_rate.toFixed(1)}%</div>
            </div>
            <div className="space-y-0.5 text-right">
              <div className="text-sm text-muted-foreground">KD Ratio</div>
              <div className="text-2xl font-bold">{stats.kd_ratio.toFixed(2)}</div>
            </div>
          </div>
          <div className="flex justify-between">
            <div className="space-y-0.5">
              <div className="text-sm text-muted-foreground">Total Matches</div>
              <div className="text-2xl font-bold">{stats.total_matches}</div>
            </div>
            <div className="space-y-0.5 text-right">
              <div className="text-sm text-muted-foreground">Recent Trends</div>
              <div className="flex space-x-2 text-sm">
                <div className="flex items-center">
                  Win/Loss {renderTrend(stats.performance_dynamics.win_loss.trend)}
                </div>
                <div className="flex items-center">
                  KDA {renderTrend(stats.performance_dynamics.kda.trend)}
                </div>
                <div className="flex items-center">
                  Rank {renderTrend(stats.performance_dynamics.rank.trend)}
                </div>
              </div>
            </div>
          </div>
        </div>

        <div>
          <h4 className="text-sm font-medium mb-2">Average Performance</h4>
          <div className="grid grid-cols-2 gap-4">
            <div className="bg-muted p-3 rounded">
              <div className="text-sm text-muted-foreground">Avg. Kills</div>
              <div className="text-xl font-medium">{stats.avg_kills_per_match.toFixed(1)}</div>
            </div>
            <div className="bg-muted p-3 rounded">
              <div className="text-sm text-muted-foreground">Avg. Deaths</div>
              <div className="text-xl font-medium">{stats.avg_deaths_per_match.toFixed(1)}</div>
            </div>
            <div className="bg-muted p-3 rounded">
              <div className="text-sm text-muted-foreground">Avg. Assists</div>
              <div className="text-xl font-medium">{stats.avg_assists_per_match.toFixed(1)}</div>
            </div>
            <div className="bg-muted p-3 rounded">
              <div className="text-sm text-muted-foreground">Avg. Match Duration</div>
              <div className="text-xl font-medium">{stats.avg_match_duration.toFixed(1)} min</div>
            </div>
          </div>
        </div>
      </CardContent>
    </Card>
  );
}; 
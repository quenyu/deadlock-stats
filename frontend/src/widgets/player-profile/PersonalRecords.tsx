import { Card, CardContent, CardHeader } from "@/shared/ui/card";

interface PersonalRecordsProps {
  records: {
    max_kills: number;
    max_assists: number;
    max_net_worth: number;
    best_kda: number;
    max_kills_match_id: string;
    max_assists_match_id: string;
    max_net_worth_match_id: string;
    best_kda_match_id: string;
  }
}

export const PersonalRecords = ({ records }: PersonalRecordsProps) => {
  if (!records) return null;
  
  return (
    <Card>
      <CardHeader>
        <h3 className="text-lg font-semibold">Personal Records</h3>
        <p className="text-sm text-muted-foreground">Best achievements in matches</p>
      </CardHeader>
      <CardContent>
        <div className="grid grid-cols-2 gap-4">
          <div className="bg-muted p-3 rounded">
            <div className="text-sm text-muted-foreground">Max Kills</div>
            <div className="text-2xl font-bold">{records.max_kills}</div>
            <div className="text-xs text-muted-foreground">Match #{records.max_kills_match_id.slice(-6)}</div>
          </div>
          <div className="bg-muted p-3 rounded">
            <div className="text-sm text-muted-foreground">Max Assists</div>
            <div className="text-2xl font-bold">{records.max_assists}</div>
            <div className="text-xs text-muted-foreground">Match #{records.max_assists_match_id.slice(-6)}</div>
          </div>
          <div className="bg-muted p-3 rounded">
            <div className="text-sm text-muted-foreground">Max Souls</div>
            <div className="text-2xl font-bold">{records.max_net_worth}</div>
            <div className="text-xs text-muted-foreground">Match #{records.max_net_worth_match_id.slice(-6)}</div>
          </div>
          <div className="bg-muted p-3 rounded">
            <div className="text-sm text-muted-foreground">Best KDA</div>
            <div className="text-2xl font-bold">{records.best_kda.toFixed(2)}</div>
            <div className="text-xs text-muted-foreground">Match #{records.best_kda_match_id.slice(-6)}</div>
          </div>
        </div>
      </CardContent>
    </Card>
  )
} 
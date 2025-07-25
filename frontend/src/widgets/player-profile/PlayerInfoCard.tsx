import { Card, CardHeader } from "@/shared/ui/card";
import { Avatar, AvatarFallback, AvatarImage } from "@/shared/ui/avatar";
import { RankBadge } from "@/shared/ui/RankBadge/RankBadge";
import { Crown } from "lucide-react";

interface PlayerInfoCardProps {
  profile: {
    nickname: string;
    avatar_url: string;
    player_rank: number;
    rank_name: string;
    rank_image: string;
    sub_rank: number;
    peak_rank?: number;
    peak_rank_name?: string;
    peak_rank_image?: string;
  };
}

export const PlayerInfoCard = ({ profile }: PlayerInfoCardProps) => {
  return (
    <Card>
      <CardHeader className="flex flex-row items-center justify-between p-4">
        <div className="flex items-center gap-4">
          <Avatar className="h-24 w-24 border-4 border-primary/10">
            <AvatarImage src={profile.avatar_url} alt={profile.nickname} />
            <AvatarFallback>{profile.nickname?.charAt(0) || "?"}</AvatarFallback>
          </Avatar>
          <div>
            <h1 className="text-3xl font-bold">{profile.nickname}</h1>
            <div className="mt-1">
              <RankBadge 
                rank={profile.rank_name}
                subrank={profile.sub_rank}
                rankImage={profile.rank_image}
              />
            </div>
          </div>
        </div>

        {profile.peak_rank && profile.peak_rank_name && (
          <div className="flex flex-col items-center gap-2 text-center">
             <div className="flex items-center gap-1 text-sm font-medium text-amber-500">
                <Crown className="h-4 w-4" />
                <span>Peak Rank</span>
              </div>
            {profile.peak_rank_image && (
              <img 
                src={profile.peak_rank_image} 
                alt={profile.peak_rank_name} 
                className="h-12 w-12" 
              />
            )}
            <span className="text-lg font-semibold">{profile.peak_rank_name}</span>
          </div>
        )}
      </CardHeader>
    </Card>
  );
}; 
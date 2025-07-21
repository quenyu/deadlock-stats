import React from 'react';

interface RankBadgeProps {
  rank: string;
  rankImage: string;
  subrank?: number;
}

const subrankSymbols = ['I', 'II', 'III', 'IV', 'V', '★'];

export const RankBadge: React.FC<RankBadgeProps> = ({ rank, subrank, rankImage }) => {
  if (!rankImage) {
    return null;
  }

  return (
    <div className="flex items-center gap-2 mt-2">
      <img src={rankImage} alt={rank} className="w-12 h-12 object-contain" />
      <div className="flex flex-col">
        <p className="text-lg font-semibold leading-none">{rank}</p>
        {subrank && (
          <p className="text-sm text-muted-foreground">Tier {subrankSymbols[subrank - 1]}</p>
        )}
      </div>
    </div>
  );
}; 
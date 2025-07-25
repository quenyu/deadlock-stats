import { Badge } from '@/shared/ui/badge'

interface RankBadgeProps {
  rank: string
  subrank?: number | null
  rankImage?: string
}

const toRoman = (num: number): string => {
  if (num < 1 || num > 6) return String(num);
  const roman: { [key: number]: string } = {
    1: 'I',
    2: 'II',
    3: 'III',
    4: 'IV',
    5: 'V',
    6: 'VI',
  };
  return roman[num];
};

export const RankBadge = ({ rank, subrank, rankImage }: RankBadgeProps) => {
  if (!rank || rank === 'Unranked') {
    return (
      <Badge variant="secondary" className="mt-1">
        Unranked
      </Badge>
    )
  }

  const subrankText = subrank ? `Tier ${toRoman(subrank)}` : ''
  const imageUrl = rankImage || ''; 

  return (
    <div className="flex items-center gap-2 mt-1">
      {imageUrl && <img src={imageUrl} alt={rank} className="h-14 w-14" />}
      <div className="flex flex-col">
        <span className="font-semibold text-lg">{rank}</span>
        {subrankText && (
          <span className="text-sm text-muted-foreground">{subrankText}</span>
        )}
      </div>
    </div>
  )
}
import { Card, CardContent, CardHeader } from "@/shared/ui/card";

interface FeaturedHeroesProps {
  featuredHeroes: {
    hero_id: number;
    hero_name: string;
    hero_image: string;
    kills?: number;
    wins?: number;
    stat_id?: number;
    stat_score?: number;
  }[]
}

export const FeaturedHeroes = ({ featuredHeroes }: FeaturedHeroesProps) => {
  if (!featuredHeroes || featuredHeroes.length === 0) return null;
  
  return (
    <Card>
      <CardHeader>
        <h3 className="text-lg font-semibold">Featured Heroes</h3>
        <p className="text-sm text-muted-foreground">Heroes showcased in player's card</p>
      </CardHeader>
      <CardContent>
        <div className="flex flex-wrap gap-4">
          {featuredHeroes.map((hero) => (
            <div key={hero.hero_id} className="flex items-center space-x-3 bg-muted p-3 rounded">
              <div className="w-10 flex-shrink-0">
                <img 
                  src={hero.hero_image} 
                  alt={hero.hero_name} 
                  className="w-full h-auto rounded-md"
                  style={{ aspectRatio: '280 / 380' }} 
                />
              </div>
              <div>
                <p className="font-medium">{hero.hero_name}</p>
                <div className="text-sm text-muted-foreground">
                  {hero.kills !== undefined && <span>{hero.kills} Kills</span>}
                  {hero.wins !== undefined && <span> • {hero.wins} Wins</span>}
                  {hero.stat_score !== undefined && <span> • Score: {hero.stat_score}</span>}
                </div>
              </div>
            </div>
          ))}
        </div>
      </CardContent>
    </Card>
  )
} 
import { PlayerSearchAdvanced } from '@/features/PlayerSearch/ui/PlayerSearchAdvanced'

export const PlayerSearchPage = () => (
  <div className="container mx-auto py-8 space-y-8">
    <div className="text-center space-y-4">
      <h1 className="text-4xl font-bold">Player Search</h1>
      <p className="text-muted-foreground text-lg max-w-2xl mx-auto">
        Discover and find players with our advanced search features. 
        Search by nickname, Steam ID, or explore popular and recently active players.
      </p>
    </div>
    <PlayerSearchAdvanced showPopular={true} showRecentlyActive={true} />
  </div>
) 
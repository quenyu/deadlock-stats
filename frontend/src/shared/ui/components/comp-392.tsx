import { UserRoundIcon } from "lucide-react"

import { Avatar, AvatarFallback } from "@/shared/ui/avatar"

export default function Component() {
  return (
    <Avatar>
      <AvatarFallback>
        <UserRoundIcon size={16} className="opacity-60" aria-hidden="true" />
      </AvatarFallback>
    </Avatar>
  )
}

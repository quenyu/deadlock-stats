"use client"

import * as React from "react"
import { Label as RadixLabel } from "@radix-ui/react-label"

import { cn } from "@/shared/lib/utils"

function Label({ className, ...props }: React.ComponentPropsWithoutRef<typeof RadixLabel>) {
  return (
    <RadixLabel
      data-slot="label"
      className={cn(
        "text-foreground text-sm leading-4 font-medium select-none group-data-[disabled=true]:pointer-events-none group-data-[disabled=true]:opacity-50 peer-disabled:cursor-not-allowed peer-disabled:opacity-50",
        className
      )}
      {...props}
    />
  )
}

export { Label }

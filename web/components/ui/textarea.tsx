import * as React from "react";

import { cn } from "@/lib/utils";

const Textarea = React.forwardRef<HTMLTextAreaElement, React.TextareaHTMLAttributes<HTMLTextAreaElement>>(
  ({ className, ...props }, ref) => (
    <textarea
      className={cn(
        "flex min-h-36 w-full rounded-md border border-ink/15 bg-white px-3 py-3 text-sm shadow-sm outline-none transition placeholder:text-ink/40 focus:border-comicRed focus:ring-2 focus:ring-comicRed/15",
        className
      )}
      ref={ref}
      {...props}
    />
  )
);
Textarea.displayName = "Textarea";

export { Textarea };

import * as React from "react";

import { cn } from "@/lib/utils";

const Input = React.forwardRef<HTMLInputElement, React.InputHTMLAttributes<HTMLInputElement>>(
  ({ className, type, ...props }, ref) => (
    <input
      type={type}
      className={cn(
        "flex h-11 w-full rounded-md border border-ink/15 bg-white px-3 text-sm shadow-sm outline-none transition placeholder:text-ink/40 focus:border-comicRed focus:ring-2 focus:ring-comicRed/15",
        className
      )}
      ref={ref}
      {...props}
    />
  )
);
Input.displayName = "Input";

export { Input };

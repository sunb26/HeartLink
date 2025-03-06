"use client";

import type { ColumnDef } from "@tanstack/react-table";
import { ArrowUpDown } from "lucide-react";
import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogTrigger,
  DialogTitle,
} from "@/components/ui/dialog";
import { CommentForm } from "@/components/ui/recordings/form";

export type Recording = {
  id: number;
  date: string;
  url: string;
};

export const columns: ColumnDef<Recording>[] = [
  {
    accessorKey: "date",
    header: ({ column }) => {
      return (
        <Button
          variant="ghost"
          onClick={() => column.toggleSorting(column.getIsSorted() === "asc")}
          className="pl-0"
        >
          Date
          <ArrowUpDown className="ml-2 h-4 w-4" />
        </Button>
      );
    },
    cell: ({ row }) => {
      return (
        <Dialog>
          <DialogTrigger asChild>
            <Button variant="ghost" className="pl-0 w-full justify-start">
              {row.original.date}
            </Button>
          </DialogTrigger>
          <DialogTitle>
            <DialogContent className="max-w-[425px] pt-10">
              <div className="bg-slate-400 p-8 rounded-lg">
                {/* biome-ignore lint/a11y/useMediaCaption: heart sounds don't require CC */}
                <audio controls>
                  <source src={row.original.url} type="audio/wav" />
                  Your browser does not support the audio element.
                </audio>
              </div>
              <div className="grid w-full gap-4">
                <CommentForm />
              </div>
            </DialogContent>
          </DialogTitle>
        </Dialog>
      );
    },
  },
];

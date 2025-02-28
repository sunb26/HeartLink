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
            <DialogContent className="sm:max-w-[425px]">
              {/* biome-ignore lint/a11y/useMediaCaption: <explanation> */}
              <audio controls autoPlay>
                <source src={row.original.url} type="audio/wav" />
                Your browser does not support the audio element.
              </audio>
            </DialogContent>
          </DialogTitle>
        </Dialog>
      );
    },
  },
];

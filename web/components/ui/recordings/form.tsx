import {
  Form,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";

import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import { z } from "zod";

import { Button } from "@/components/ui/button";
import { Textarea } from "@/components/ui/textarea";
import { useState } from "react";

const FormSchema = z.object({
  comments: z.string().min(2, {
    message: "Comments must be at least 2 characters.",
  }),
});

export function CommentForm() {
  const [submitted, setSubmitted] = useState(false);

  const form = useForm<z.infer<typeof FormSchema>>({
    resolver: zodResolver(FormSchema),
  });

  function onSubmit(data: z.infer<typeof FormSchema>) {
    const reqOptioons = {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(data),
    };
    fetch("https://heartlink.free.beeceptor.com/todos", reqOptioons).then(
      (res) => {
        if (res.ok) {
          setSubmitted(true);
        } else {
          form.setError("comments", {
            type: "server",
            message: "Something went wrong. Please try again.",
          });
        }
      }
    );
  }

  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-6">
        <FormField
          control={form.control}
          name="comments"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Comments</FormLabel>
              <FormControl>
                <Textarea
                  placeholder="Type your comments here..."
                  className="md:h-[425]"
                  {...field}
                />
              </FormControl>
              <FormDescription>
                Provide comments on this recording for the patient.
              </FormDescription>
              <FormMessage />
            </FormItem>
          )}
        />
        <Button
          type="submit"
          className="w-full text-white bg-gradient-to-r from-red-400 via-red-500 to-red-600 hover:bg-gradient-to-br focus:ring-4 focus:outline-none focus:ring-red-300 dark:focus:ring-red-800 shadow-lg shadow-red-500/50 dark:shadow-lg dark:shadow-red-800/80 font-medium rounded-lg text-sm px-10 py-4 text-center me-2 mb-2"
        >
          Submit
        </Button>
        {submitted && (
          <FormMessage className="text-green-900 font-bold">
            Your comments have been submitted successfully.
          </FormMessage>
        )}
      </form>
    </Form>
  );
}

"use client";

import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import { z } from "zod";
import { Button } from "@/components/ui/button";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { useState } from "react";

type reqBody = {
  physicianId: string,
  firstName: string,
  lastName: string,
  email: string,
  dob: string,
  sex: string,
  height: number,
  weight: number,
}


const formSchema = z.object({
  firstName: z.string().min(2, {
    message: "First Name must be at least 2 characters.",
  }),
  lastName: z.string().min(2, {
    message: "Last Name must be at least 2 characters.",
  }),
  email: z.string().email({
    message: "Please enter a valid email address.",
  }),
  dob: z.coerce.date({
    required_error: "Date of birth is required",
    invalid_type_error: "Invalid date format",
  }),
  sex: z.enum(["M", "F", "NB", ""], {
    message: "Please enter capital 'M', 'F', or 'NB'.",
  }),
  height: z.coerce.number().min(0, {
    message: "Please enter a valid height.",
  }),
  weight: z.coerce.number().min(0, {
    message: "Please enter a valid weight.",
  }),
});

export function RegisterPatientForm( { physicianId }: { physicianId?: string }) {
  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      firstName: "",
      lastName: "",
      email: "",
      sex: "",
      height: 0,
      weight: 0,
    },
  });
  const [submitted, setSubmitted] = useState(false);

  function onSubmit(values: z.infer<typeof formSchema>) {
    if (physicianId == null) {
      form.setError("weight", {
        type: "weight",
        message: "Physician ID not found",
      });
      return;
    }
    const rb : reqBody = {
      physicianId: physicianId,
      firstName: values.firstName,
      lastName: values.lastName,
      email: values.email,
      sex: values.sex,
      dob: values.dob.toISOString(),
      height: values.height,
      weight: values.weight,
    };
    const reqOptions = {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(rb),
    };
    fetch(`${process.env.NEXT_PUBLIC_SERVER_URL}/CreatePatient`, reqOptions).then(
      (res) => {
        if (res.ok) {
          setSubmitted(true);
          console.log("Patient created successfully");
        } else {
          console.log("Error creating patient");
          console.log(res);
          form.setError("weight", {
            type: "server",
            message: "Something went wrong. Please try again.",
          });
        }
      }
    );
  }

  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-3 max-h-3/4">
        <FormField
          control={form.control}
          name="firstName"
          render={({ field }) => (
            <FormItem>
              <FormLabel>First Name</FormLabel>
              <FormControl>
                <Input placeholder="First Name" {...field} />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />
        <FormField
          control={form.control}
          name="lastName"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Last Name</FormLabel>
              <FormControl>
                <Input placeholder="Last Name" {...field} />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />
        <FormField
          control={form.control}
          name="email"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Email</FormLabel>
              <FormControl>
                <Input placeholder="email@example.com" {...field} />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />
        <FormField
          control={form.control}
          name="dob"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Date of Birth</FormLabel>
              <FormControl>
                <Input
                  type="date"
                  placeholder="YYYY-MM-DD"
                  value={
                    field.value
                      ? new Date(field.value).toISOString().split("T")[0]
                      : ""
                  }
                  onChange={(e) => field.onChange(new Date(e.target.value))}
                />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />
        <FormField
          control={form.control}
          name="sex"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Sex</FormLabel>
              <FormControl>
                <Input placeholder="M or F or NB" {...field} />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />
        <FormField
          control={form.control}
          name="height"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Height (cm)</FormLabel>
              <FormControl>
                <Input placeholder="Height" {...field} />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />
        <FormField
          control={form.control}
          name="weight"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Weight (lbs)</FormLabel>
              <FormControl>
                <Input placeholder="Weight" {...field} />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />
        <Button type="submit">Confirm Registration</Button>
      </form>
      {submitted && (
          <FormMessage className="text-green-900 font-bold">
            Patient created successfully.
          </FormMessage>
        )}
    </Form>
  );
}

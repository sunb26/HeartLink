import { type Recording, columns } from "@/components/ui/recordings/columns";
import { DataTable } from "@/components/ui/recordings/data-table";
import Image from "next/image";
import React from "react";

type Patient = {
  firstName: string;
  lastName: string;
  email: string;
  age: number;
  sex: string;
  weight: string;
  height: string;
  recordings: Recording[];
};

async function getPatient(): Promise<Patient> {
  return {
    firstName: "John",
    lastName: "Doe",
    email: "example@gmail.com",
    age: 25,
    sex: "M",
    weight: "160 lbs",
    height: "180 cm",
    recordings: [
      {
        id: 1,
        date: "2024-02-21T14:30:00Z",
        url: "https://firebasestorage.googleapis.com/v0/b/heartlink-6fee0.firebasestorage.app/o/recordings%2Fchillin39-20915.wav?alt=media&token=7ac4ed0a-d427-47cc-9172-ca727c596179",
      },
      {
        id: 2,
        date: "2024-02-20T09:15:00Z",
        url: "https://firebasestorage.googleapis.com/v0/b/heartlink-6fee0.firebasestorage.app/o/recordings%2Fchillin39-20915.wav?alt=media&token=7ac4ed0a-d427-47cc-9172-ca727c596179",
      },
      {
        id: 3,
        date: "2024-02-19T16:45:00Z",
        url: "https://firebasestorage.googleapis.com/v0/b/heartlink-6fee0.firebasestorage.app/o/recordings%2Fchillin39-20915.wav?alt=media&token=7ac4ed0a-d427-47cc-9172-ca727c596179",
      },
      {
        id: 4,
        date: "2024-02-18T11:20:00Z",
        url: "https://firebasestorage.googleapis.com/v0/b/heartlink-6fee0.firebasestorage.app/o/recordings%2Fchillin39-20915.wav?alt=media&token=7ac4ed0a-d427-47cc-9172-ca727c596179",
      },
      {
        id: 5,
        date: "2024-02-17T13:00:00Z",
        url: "https://firebasestorage.googleapis.com/v0/b/heartlink-6fee0.firebasestorage.app/o/recordings%2Fchillin39-20915.wav?alt=media&token=7ac4ed0a-d427-47cc-9172-ca727c596179",
      },
      {
        id: 6,
        date: "2024-02-16T15:30:00Z",
        url: "https://firebasestorage.googleapis.com/v0/b/heartlink-6fee0.firebasestorage.app/o/recordings%2Fchillin39-20915.wav?alt=media&token=7ac4ed0a-d427-47cc-9172-ca727c596179",
      },
      {
        id: 7,
        date: "2024-02-15T10:45:00Z",
        url: "https://firebasestorage.googleapis.com/v0/b/heartlink-6fee0.firebasestorage.app/o/recordings%2Fchillin39-20915.wav?alt=media&token=7ac4ed0a-d427-47cc-9172-ca727c596179",
      },
      {
        id: 8,
        date: "2024-02-14T12:15:00Z",
        url: "https://firebasestorage.googleapis.com/v0/b/heartlink-6fee0.firebasestorage.app/o/recordings%2Fchillin39-20915.wav?alt=media&token=7ac4ed0a-d427-47cc-9172-ca727c596179",
      },
      {
        id: 9,
        date: "2024-02-13T09:00:00Z",
        url: "https://firebasestorage.googleapis.com/v0/b/heartlink-6fee0.firebasestorage.app/o/recordings%2Fchillin39-20915.wav?alt=media&token=7ac4ed0a-d427-47cc-9172-ca727c596179",
      },
      {
        id: 10,
        date: "2024-02-12T14:20:00Z",
        url: "https://firebasestorage.googleapis.com/v0/b/heartlink-6fee0.firebasestorage.app/o/recordings%2Fchillin39-20915.wav?alt=media&token=7ac4ed0a-d427-47cc-9172-ca727c596179",
      },
      {
        id: 11,
        date: "2024-02-11T16:30:00Z",
        url: "https://firebasestorage.googleapis.com/v0/b/heartlink-6fee0.firebasestorage.app/o/recordings%2Fchillin39-20915.wav?alt=media&token=7ac4ed0a-d427-47cc-9172-ca727c596179",
      },
      {
        id: 12,
        date: "2024-02-10T11:45:00Z",
        url: "https://firebasestorage.googleapis.com/v0/b/heartlink-6fee0.firebasestorage.app/o/recordings%2Fchillin39-20915.wav?alt=media&token=7ac4ed0a-d427-47cc-9172-ca727c596179",
      },
      {
        id: 13,
        date: "2024-02-09T10:00:00Z",
        url: "https://firebasestorage.googleapis.com/v0/b/heartlink-6fee0.firebasestorage.app/o/recordings%2Fchillin39-20915.wav?alt=media&token=7ac4ed0a-d427-47cc-9172-ca727c596179",
      },
      // ...
    ],
  };
}

export default async function PatientPage() {
  const patient = await getPatient();
  return (
    <div>
      <div className="bg-off-white">
        <div className="container mx-auto flex items-center gap-16">
          <div className="w-48 h-48 rounded-full overflow-hidden">
            <Image
              src={"/fallback-avatar.svg"}
              width={200}
              height={200}
              alt="Patient Avatar"
              className="w-full h-full object-cover"
            />
          </div>
          <p className="text-3xl text-left py-10 font-[Syne] whitespace-pre-line">
            First Name: {patient.firstName}
            <br />
            Last Name: {patient.lastName}
            <br />
            Email: {patient.email}
            <br />
            Age: {patient.age}
            <br />
            Sex: {patient.sex}
            <br />
            Weight: {patient.weight}
            <br />
            Height: {patient.height}
          </p>
        </div>
      </div>
      <div className="container mx-auto py-10">
        <h2 className="text-2xl font-bold pb-4 font-[Syne]">Recordings</h2>
        <DataTable columns={columns} data={patient.recordings} />
      </div>
    </div>
  );
}

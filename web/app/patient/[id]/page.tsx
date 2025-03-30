import { columns } from "@/components/ui/recordings/columns";
import { DataTable } from "@/components/ui/recordings/data-table";
import Image from "next/image";
import React from "react";

export default async function PatientPage({params}: {params: Promise<{ id: string }>}) {
  const { id } = await params;
  console.log("PatientId:", id);
  const response = await fetch(`${process.env.NEXT_PUBLIC_SERVER_URL}/GetPatient?patientId=${id}`, {
    method: "GET",
    headers: {
      "Content-Type": "application/json",
    },
  }).catch((error) => {
    console.error("Error:", error);
  });
  if (!response) {
    return [];
  }
  const patient = await response.json();
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
            Weight: {patient.weight} lbs
            <br />
            Height: {patient.height} cm
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

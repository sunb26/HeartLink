import { columns } from "@/components/ui/patients/columns";
import { DataTable } from "@/components/ui/patients/data-table";
import { currentUser } from "@clerk/nextjs/server";


export default async function PhysicianPage() {
  const physician = await currentUser();
  const response = await fetch(`${process.env.NEXT_PUBLIC_SERVER_URL}/ListPatients?physicianid=${physician?.id}`, {
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
  const body = await response.json();
  const data = body.patients;

  return (
    <div>
      <div className="bg-off-white">
        <h1 className="text-4xl font-bold text-center py-10 font-[Syne]">
          Welcome Dr. {physician?.fullName}
        </h1>
      </div>
      <div className="container mx-auto py-10">
        <h2 className="text-2xl font-bold pb-4 font-[Syne]">Your Patients</h2>
        <DataTable columns={columns} data={data} />
      </div>
    </div>
  );
}
